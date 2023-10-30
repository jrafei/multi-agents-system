package restserveragent

import (
	"errors"
	"fmt"
	rad_t "ia04/agt"
	"ia04/comsoc"
	com "ia04/comsoc"
	"net/http"
	"sync"
	"time"
)

type RestBallotAgent struct {
	sync.Mutex
	id          string
	rule        string
	deadline    string
	voter_ids   map[string]bool // Le booléen permet de savoir si l'agent a déjà voté
	nb_alts     int
	tiebreak    []com.Alternative
	server_chan chan rad_t.RequestVoteBallot // canal de communication entre le scrutin et le serveur (requetes de type vote ou result)
	profile     com.Profile
	options     [][]int
}

func NewRestBallotAgent(i string, ru string, d string, vot_ids []string, alts int, tieb []com.Alternative, ch chan rad_t.RequestVoteBallot) *RestBallotAgent {
	voters := make(map[string]bool, 0)
	for _, id := range vot_ids {
		// On signale qu'aucun voteur n'a encore voté
		voters[id] = false
	}
	return &RestBallotAgent{id: i, rule: ru, deadline: d, voter_ids: voters, nb_alts: alts, tiebreak: tieb, server_chan: ch, profile: make(com.Profile, 0), options: make([][]int, 0)}
}

/*
traite les requetes que le channel 'server_chan' a entendu
*/
func (rsa *RestBallotAgent) Start(chan rad_t.RequestVoteBallot) {
	// si le channel reçoit une demande, on lance la méthode associée
	for {
		var resp rad_t.RequestVoteBallot
		req := <-rsa.server_chan
		// Selection de l'action à effectuer
		switch req.Action {
		case "vote":
			resp = rsa.Vote(req)
		case "result":
			resp = rsa.Result()
		default:
			resp.StatusCode = http.StatusBadRequest
			resp.Msg = "bad request, unknown process for ballot"
		}
		// Transmission de la réponse au serveur
		rsa.server_chan <- resp

	}
}

/*
ajout d'un vote s'il est valide
*/
func (rba *RestBallotAgent) Vote(vote rad_t.RequestVoteBallot) (resp rad_t.RequestVoteBallot) {
	rba.Lock()
	defer rba.Unlock()
	// Vérification de la deadline
	if rba.deadline <= time.Now().Format(time.RFC3339) {
		resp.StatusCode = http.StatusGatewayTimeout
		resp.Msg = "la deadline est dépassée"
		return
	}

	// Vérification de l'AgentID
	has_voted, exists := rba.voter_ids[vote.AgentID]
	if !exists {
		resp.StatusCode = http.StatusBadRequest
		resp.Msg = "bad request, le voteur n'est pas sur la liste"
		return
	}
	if has_voted {
		resp.StatusCode = http.StatusForbidden
		resp.Msg = "vote déjà effectué"
		return
	}

	// Vérification des préférences
	prefs := make([]comsoc.Alternative, len(vote.Preferences))
	alts := make([]comsoc.Alternative, rba.nb_alts)
	for i, _ := range alts {
		alts[i] = comsoc.Alternative(i + 1)
	}
	for i, _ := range prefs {
		prefs[i] = comsoc.Alternative(vote.Preferences[i])
	}

	if (comsoc.CheckProfile(prefs, alts) != nil) || (len(alts) != len(prefs)) { // TODO : vérifier si la 2ème condition doit être intégrer cette vérification dans checkProfil()
		resp.StatusCode = http.StatusBadRequest
		resp.Msg = " [Agent " + vote.RequestVote.AgentID + "] bad request, les préférences ne sont pas conformes"
		return
	}

	// Ajout des options
	if vote.Options != nil && len(vote.Options) > 0 && (vote.Options[0] <= rba.nb_alts && vote.Options[0] >= 1) { // On part du principe que la première valeur est un seuil de vote (cf.Approval)
		rba.options = append(rba.options, vote.Options)
	} else if rba.rule == "approval" {
		// si pas de seuil de préférence pour la méthode par approbation, erreur !
		resp.StatusCode = http.StatusBadRequest
		resp.Msg = "bad request, aucun seuil de préférence saisi"
		return
	}

	// Ajout des préférences dans le profil
	rba.profile = append(rba.profile, prefs)

	rba.voter_ids[vote.AgentID] = true // on indique que l'agent a voté
	resp.StatusCode = http.StatusOK
	resp.Msg = "vote pris en compte"

	/********DEBUG********/
	fmt.Println("-----------------")
	fmt.Printf("[DBG] [%s] Updated ballot after /vote : \n", vote.AgentID)
	fmt.Println(rba.id)
	fmt.Println(rba.deadline)
	fmt.Println(rba.nb_alts)
	fmt.Println(rba.rule)
	fmt.Println(rba.profile)
	fmt.Println(rba.options)
	fmt.Println("-----------------")
	/*********************/

	/********DEBUG********/
	fmt.Println("-----------------")
	fmt.Printf("[DBG] [%s] Response /vote from ballot to server : \n", vote.AgentID)
	fmt.Println(rba.id)
	fmt.Println(rba.deadline)
	fmt.Println(rba.nb_alts)
	fmt.Println(rba.rule)
	fmt.Println(rba.profile)
	fmt.Println(rba.options)
	fmt.Println("-----------------")
	/*********************/

	return
}

/*
 */
func (rsa *RestBallotAgent) Result() (resp rad_t.RequestVoteBallot) {
	rsa.Lock()
	defer rsa.Unlock()
	// Vérification de la deadline
	if rsa.deadline > time.Now().Format(time.RFC3339) {
		resp.StatusCode = http.StatusTooEarly
		resp.Msg = "too early"
		return
	}
	var ranking []comsoc.Alternative
	var err error
	switch rsa.rule {
	case "majority":
		ranking, err = comsoc.SWFFactory(com.MajoritySWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile)
	case "borda":
		ranking, err = comsoc.SWFFactory(com.BordaSWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile)
	case "approval":
		// récupération des seuils pour le vote
		thresholds := make([]int, len(rsa.voter_ids))
		for i, _ := range rsa.options {
			thresholds[i] = rsa.options[i][0] // On part du principe que c'est la première valeur
		}
		ranking, err = comsoc.SWFFactoryOptions[int](com.ApprovalSWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile, thresholds)
	case "stv":
		ranking, err = comsoc.SWFFactoryOptions[comsoc.Alternative](com.STV_SWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile, rsa.tiebreak)
	case "copeland":
		ranking, err = comsoc.SWFFactory(com.CopelandSWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile)
	case "condorcet":
		ranking, err = comsoc.CondorcetWinner(rsa.profile)
		fmt.Println(ranking, err)
	default :
		err = errors.New("unknown rule")
	}


	if err == nil {
		resp.StatusCode = http.StatusOK
		if len(ranking) > 1 {
			resp.Ranking = make([]int, rsa.nb_alts)
			for i, _ := range ranking {
				resp.Ranking[i] = int(ranking[i])
			}
			resp.Winner = resp.Ranking[0]
		} else if len(ranking) == 1 {
			resp.Winner = int(ranking[0])
		}

	} else {
		resp.StatusCode = http.StatusInternalServerError
		resp.Msg = "internal server error, " + err.Error()
	}

	return
}
