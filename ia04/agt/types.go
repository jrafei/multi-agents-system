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
	Result int `json:"res"`
}
