package comsoc

import "errors"

/*
 Reference du méthode : https://en.wikipedia.org/wiki/Kemeny%E2%80%93Young_method
*/

func Kemeny(p Profile, orderedAlts []Alternative) ([]Alternative, error) {

	if len(p) == 0 {
		return nil, errors.New("profil is empty")
	}

	err := CheckProfileAlternative(p, p[0])
	if err != nil {
		return nil, err
	}

	battle := CountIsPref(p) // de type map[AltTuple] int -> le nombre de fois où 'a' bat 'b'

	//calcule toutes les possibilités de ranking  -> permutation d'une préférence
	rankings := [][]Alternative{}
	permute(p[0], 0, &rankings)

	bestRank := [][]Alternative{rankings[0]}
	bestScore := calculateScore(bestRank[0], battle)

	for i := 1; i < len(rankings); i++ {
		r := rankings[i]
		s := calculateScore(r, battle)
		if s > bestScore {
			bestRank = make([][]Alternative, 0)
			bestRank = append(bestRank, r)
			bestScore = s
		}
		if s == bestScore {
			bestRank = append(bestRank, r)
		}
	}

	var bestWinners []Alternative

	if len(bestRank) > 1 {
		// application de tiebreak
		for _, r := range bestRank {
			bestWinners = append(bestWinners, r[0])
		}
		return []Alternative{meilleurElement(bestWinners, orderedAlts)}, nil
	}

	//return []Alternative{bestRank[0][0]}, nil
	return bestRank[0], nil
}

/*
	Retourne le score du classement
*/
func calculateScore(ranking []Alternative, battle map[AltTuple]int) int {
	res := 0
	for x, _ := range ranking {
		for y := x + 1; y < len(ranking); y++ {
			res += battle[AltTuple{Alternative(ranking[x]), Alternative(ranking[y])}]
		}
	}
	return res

}
