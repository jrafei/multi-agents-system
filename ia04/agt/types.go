package agt

import (
	coms "ia04/comsoc"
)

// Requête pour la création d'un agent (transfert via requête http)
type RequestBallot struct {
	Rule     string   `json:"rule"`
	Deadline string   `json:"deadline"`
	Voters   []string `json:"voter-ids"`
	Nb_alts  int      `json:"#alts"`
	Tiebreak []int    `json:"tie-break"`
}

// Requête pour la prise en compte d'un vote, et le résultat d'un scrutin (transfert via requête http)
type RequestVote struct {
	AgentID     string             `json:"agent-id,omitempty"`
	BallotID    string             `json:"ballot-id"`
	Preferences []coms.Alternative `json:"prefs,omitempty"`
	Options     []int              `json:"options,omitempty"`
}

// Requête échangée entre le ballot et le serveur (requete interne)
type RequestVoteBallot struct {
	*RequestVote        //renseigné par le serveur
	Action       string //renseigné par le serveur, ex : vote, result...
	StatusCode   int    //renseigné par le ballot
	Msg          string //renseigné par le ballot
	Winner       int    //renseigné par le ballot
	Ranking      []int  //renseigné par le ballot
}

// Requête de réponse générale (transfert via requête http)
type Response struct {
	Ballot_id string `json:"ballot-id,omitempty"`
	Winner    int    `json:"winner,omitempty"`
	Ranking   []int  `json:"ranking,omitempty"`
	Status    int    `json:"status,omitempty"`
}
