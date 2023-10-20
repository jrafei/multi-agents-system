package agt

// Requête pour la création d'un agent
type RequestBallot struct {
	Rule     string   `json:"rule"`
	Deadline string   `json:"deadline"`
	Voters   []string `json:"voter-ids"`
	Nb_alts  int      `json:"#alts"`
	Tiebreak []int    `json:"tie-break"`
}

// Requête pour la prise en compte d'un vote, et le résultat d'un scrutin
type RequestVote struct {
	AgentID string `json:"agent-id,omitempty"`
	BallotID string `json:"ballot-id"`
	Preferences []int `json:"prefs,omitempty"`
	Options []int`json:"options,omitempty"`
}

// Requête échangée entre le ballot et le serveur
type RequestVoteBallot struct {
	*RequestVote
	Action string
	StatusCode int
	Msg string
}

// Requête de réponse générale
type Response struct {
	Ballot_id string `json:"ballot-id"`
	Winner int  `json:"winner,omitempty"`
	Ranking []int `json:"ranking,omitempty"`
}
