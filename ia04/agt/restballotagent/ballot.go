package restserveragent

import (
	"fmt"
	rad_t "ia04/agt"
	"ia04/comsoc"
	com "ia04/comsoc"
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
	server_chan chan rad_t.RequestVoteBallot
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

func (rsa *RestBallotAgent) Start(chan rad_t.RequestVoteBallot) {
	// si le channel reçoit une demande, on lace la méthode associée
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
			resp.StatusCode = 400
			resp.Msg = "bad request, unknown process for ballot"
		}
		// Transmission de la réponse au serveur
		rsa.server_chan <- resp

	}
}

func (rsa *RestBallotAgent) Vote(vote rad_t.RequestVoteBallot) (resp rad_t.RequestVoteBallot) {

	// Vérification de la deadline
	if rsa.deadline <= time.Now().Format(time.RFC3339) {
		resp.StatusCode = 503
		resp.Msg = "la deadline est dépassée"
		return
	}

	// Vérification de l'AgentID
	has_voted, exists := rsa.voter_ids[vote.AgentID]
	if !exists {
		resp.StatusCode = 400
		resp.Msg = "bad request, le voteur n'est pas sur la liste"
		return
	}
	if has_voted {
		resp.StatusCode = 403
		resp.Msg = "vote déjà effectué"
		return
	}

	// Vérification des préférences
	prefs := make([]comsoc.Alternative, len(vote.Preferences))
	alts := make([]comsoc.Alternative, rsa.nb_alts)
	for i, _ := range alts {
		alts[i] = comsoc.Alternative(i + 1)
	}
	for i, _ := range prefs {
		prefs[i] = comsoc.Alternative(vote.Preferences[i])
	}

	if comsoc.CheckProfile(prefs, alts) != nil {
		resp.StatusCode = 400
		resp.Msg = "bad request, les préférences ne sont pas conformes"
		return
	}

	// Ajout des préférences dans le profil
	rsa.profile = append(rsa.profile, prefs)

	// Ajout des options
	if vote.Options != nil && len(vote.Options) > 0 && (vote.Options[0] <= rsa.nb_alts && vote.Options[0] >= 1) { // On part du principe que la première valeur est un seuil de vote (cf.Approval)
		rsa.options = append(rsa.options, vote.Options)
	} else if rsa.rule == "approval" {
		// si pas de seuil de préférence pour la méthode par approbation, erreur !
		resp.StatusCode = 400
		resp.Msg = "bad request, aucun seuil de préférence saisi"
		return
	}

	rsa.voter_ids[vote.AgentID] = true // on indique que l'agent a voté
	resp.StatusCode = 200
	resp.Msg = "vote pris en compte"

	/********DEBUG********/
	fmt.Println("-----------------")
	fmt.Println("[DBG] Updated ballot after /vote :")
	fmt.Println(rsa.id)
	fmt.Println(rsa.deadline)
	fmt.Println(rsa.nb_alts)
	fmt.Println(rsa.rule)
	fmt.Println(rsa.profile)
	fmt.Println(rsa.options)
	fmt.Println("-----------------")
	/*********************/

	/********DEBUG********/
	fmt.Println("-----------------")
	fmt.Println("[DBG] Response /vote from ballot to server :")
	fmt.Println(rsa.id)
	fmt.Println(rsa.deadline)
	fmt.Println(rsa.nb_alts)
	fmt.Println(rsa.rule)
	fmt.Println(rsa.profile)
	fmt.Println(rsa.options)
	fmt.Println("-----------------")
	/*********************/

	return
}

func (rsa *RestBallotAgent) Result() (resp rad_t.RequestVoteBallot) {
	// Vérification de la deadline
	if rsa.deadline > time.Now().Format(time.RFC3339) {
		resp.StatusCode = 425
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
		// TODO : test
		// récupération des seuils de vote
		thresholds := make([]int, len(rsa.voter_ids))
		for i, _ := range rsa.options {
			thresholds[i] = rsa.options[i][0] // On part du principe que c'est la première valeur
		}
		ranking, err = comsoc.SWFFactoryOptions[int](com.ApprovalSWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile, thresholds)
	case "stv":
		// TODO : test
		ranking, err = comsoc.SWFFactoryOptions[comsoc.Alternative](com.STV_SWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile, rsa.tiebreak)
	case "copeland":
		// TODO : test
		ranking, err = comsoc.SWFFactory(com.CopelandSWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile)

	}
	if err == nil {
		resp.StatusCode = 200
		resp.Msg = "OK"
		resp.Ranking = make([]int, rsa.nb_alts)
		for i, _ := range ranking {
			resp.Ranking[i] = int(ranking[i])
		}
		resp.Winner = resp.Ranking[0]
	} else {
		resp.StatusCode = 500
		resp.Msg = "internal server error, " + err.Error()
	}
	return
}
