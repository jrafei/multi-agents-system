package restserveragent

import (
	"net/http"
	"sync"
	"time"

	utils "ia04/agt/utils"
	comsoc "ia04/comsoc"
)

type RestBallotAgent struct {
	sync.Mutex
	id          string
	rule        string
	deadline    string
	voter_ids   map[string]bool // Le booléen permet de savoir si l'agent a déjà voté
	nb_alts     int
	tiebreak    []comsoc.Alternative
	server_chan chan utils.RequestVoteBallot // canal de communication entre le scrutin et le serveur (requetes de type vote ou result)
	profile     comsoc.Profile
	options     [][]int
}

/*
======================================

	  @brief :
	  'Constructeur de la classe.'
	  @params :
		- 'id' : identifiant unique du ballot
		- 'rule' : méthode de vote
		- 'deadline' : date de fin de prise en compte des votes
		- 'voters_list' : liste des identifiants des voteurs
		- 'nb_alts' : nombre d alternatives pour le vote
		- 'tieb' : liste du classement des alternatives pour TieBreak
		- 'server_chan' : channel pour échange ballot-server
	  @returned :
	    -  Un pointeur sur le ballot créé.

======================================
*/
func NewRestBallotAgent(id string, rule string, deadline string, voters_list []string, nb_alts int, tieb []comsoc.Alternative, server_chan chan utils.RequestVoteBallot) *RestBallotAgent {
	voters := make(map[string]bool, 0)
	for _, id := range voters_list {
		// On signale qu'aucun voteur n'a encore voté
		voters[id] = false
	}
	return &RestBallotAgent{id: id, rule: rule, deadline: deadline, voter_ids: voters, nb_alts: nb_alts, tiebreak: tieb, server_chan: server_chan, profile: make(comsoc.Profile, 0), options: make([][]int, 0)}
}

/*
======================================

	@brief:
	'Procédure de mise en fonction du ballot. Elle écoute et traite les requêtes échangées avec le serveur.'

======================================
*/
func (rsa *RestBallotAgent) Start() {
	for {
		var resp utils.RequestVoteBallot
		req := <-rsa.server_chan
		// Selection de l'action à effectuer
		switch req.Action {
		case "vote":
			resp = rsa.vote(req)
		case "result":
			resp = rsa.result()
		default:
			resp.StatusCode = http.StatusBadRequest
			resp.Msg = "Action inconnue."
		}
		// Transmission de la réponse au serveur
		rsa.server_chan <- resp
	}
}

/*
======================================

	  @brief :
	  'Méthode pour la prise en compte d un vote.'
	  @params :
		- 'vote' : requête entrante de type RequestVoteBallot
	  @returned :
	    - 'resp' : requête sortante (réponse) de type RequestVoteBallot

======================================
*/
func (rba *RestBallotAgent) vote(vote utils.RequestVoteBallot) (resp utils.RequestVoteBallot) {
	rba.Lock()
	defer rba.Unlock()
	// Vérification de la deadline
	if rba.deadline <= time.Now().Format(time.RFC3339) {
		resp.StatusCode = http.StatusGatewayTimeout
		resp.Msg = "La deadline est dépassée."
		return
	}

	// Vérification de l'AgentID
	has_voted, exists := rba.voter_ids[vote.AgentID]
	if !exists {
		resp.StatusCode = http.StatusBadRequest
		resp.Msg = "Le voteur n'est pas sur la liste."
		return
	}
	if has_voted {
		resp.StatusCode = http.StatusForbidden
		resp.Msg = "Vote déjà effectué."
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
		resp.Msg = " Les préférences ne sont pas conformes."
		return
	}
	
	// Ajout des options
	if vote.Options != nil && len(vote.Options) > 0 && (vote.Options[0] <= rba.nb_alts && vote.Options[0] >= 1) { // On part du principe que la première valeur est un seuil de vote (cf.Approval)
		rba.options = append(rba.options, vote.Options)
	} else if rba.rule == "approval" {
		// si pas de seuil de préférence pour la méthode par approbation, erreur !
		resp.StatusCode = http.StatusBadRequest
		resp.Msg = "Aucun seuil de préférence saisi."
		return
	}

	// Ajout des préférences dans le profil
	rba.profile = append(rba.profile, prefs)

	rba.voter_ids[vote.AgentID] = true // on indique que l'agent a voté
	resp.StatusCode = http.StatusOK
	resp.Msg = "Vote pris en compte."

	return
}

/*
======================================

	@brief :
	'Méthode pour l obtention du résultat du vote'
	@returned :
	   - 'resp' : requête sortante (réponse) de type RequestVoteBallot

======================================
*/
func (rsa *RestBallotAgent) result() (resp utils.RequestVoteBallot) {
	rsa.Lock()
	defer rsa.Unlock()
	// Vérification de la deadline
	if rsa.deadline > time.Now().Format(time.RFC3339) {
		resp.StatusCode = http.StatusTooEarly
		resp.Msg = "Le vote n'est pas encore clôturé."
		return
	}
	var ranking []comsoc.Alternative
	var err error
	switch rsa.rule {
	case "majority":
		ranking, err = comsoc.SWFFactory(comsoc.MajoritySWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile)
	case "borda":
		ranking, err = comsoc.SWFFactory(comsoc.BordaSWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile)
	case "approval":
		// récupération des seuils pour le vote
		thresholds := make([]int, len(rsa.voter_ids))
		for i, _ := range rsa.options {
			thresholds[i] = rsa.options[i][0] // On part du principe que c'est la première valeur
		}
		ranking, err = comsoc.SWFFactoryOptions[int](comsoc.ApprovalSWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile, thresholds)
	case "stv":
		ranking, err = comsoc.SWFFactoryOptions[comsoc.Alternative](comsoc.STV_SWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile, rsa.tiebreak)
	case "copeland":
		ranking, err = comsoc.SWFFactory(comsoc.CopelandSWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile)
	case "condorcet":
		ranking, err = comsoc.CondorcetWinner(rsa.profile)
	case "young-condorcet":
		ranking, err = comsoc.YoungCondorcet(rsa.profile,rsa.tiebreak)
	case "kramer-simpson":
		ranking, err = comsoc.SWFFactory(comsoc.KramerSimpson_SWF, comsoc.TieBreakFactory(rsa.tiebreak))(rsa.profile)
	default:
		resp.StatusCode = http.StatusNotImplemented
		resp.Msg = "Méthode de vote non implémentée."
		return
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
		} else if len(ranking) == 0 {
			resp.Msg = "Aucun gagnant..."
		}

	} else {
		resp.StatusCode = http.StatusInternalServerError
		resp.Msg = "Erreur interne : " + err.Error()
	}

	return
}
