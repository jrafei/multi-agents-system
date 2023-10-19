package agt

type RequestBallot struct {
	Rule     string   `json:"rule"`
	Deadline string   `json:"deadline"`
	Voters   []string `json:"voter-ids"`
	Nb_alts  int      `json:"#alts"`
	Tiebreak []int    `json:"tie-break"`
}

type RequestVote struct {
}

type Response struct {
	Ballot_id string `json:"ballot-id"`
	Winner int  `json:"winner,omitempty"`
	Ranking []int `json:"ranking,omitempty"`
}
