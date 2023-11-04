package comsoc

import "errors"

/*
Kemeny : Calculer la relation de préférence la plus proche du profil
(pour la somme pour chaque votant du nombre de paires de candidats qui
ne sont pas dans le même ordre). Le meilleur candidat est celui qui est
rangé premier pour cette relation.
*/

func Kemeny(p Profile, orderedRankings [][]Alternative) (Alternative, error) {

	if len(p) == 0 {
		return -1, errors.New("profil is empty")
	}

	err := CheckProfileAlternative(p, p[0])
	if err != nil {
		return -1, err
	}

	battle := CountIsPref(p) // de type map[AltTuple] int -> le nombre de fois où 'a' bat 'b'

	//calcule toutes les possibilités de  -> permutation d'une préférence
	rankings := [][]Alternative{}
	permute(p[0], 0, &rankings)

	bestRank := [][]Alternative{rankings[0]}
	bestScore := calculateScore(bestRank[0], battle)

	for i := 1; i < len(rankings); i++ {
		r := rankings[i]
		s := calculateScore(r, battle)
		if s > bestScore {
			bestRank := make([][]Alternative, 0)
			bestRank = append(bestRank, r)
			bestScore = s
		}
		if s == bestScore {
			bestRank = append(bestRank, r)
		}
	}

	var best []Alternative
	if len(bestRank) > 1 {
		// application de tiebreak
		best = tiebreak_kemeny(bestRank, orderedRankings)
	} else {
		best = bestRank[0]
	}
	return best[0], nil
}

/*
	Retourne le score du rang
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

func tiebreak_kemeny(winners [][]Alternative, orderedList [][]Alternative) []Alternative {
	if len(orderedList) == 0 {
		// on prend le premier
		return winners[0]
	}

	best := winners[0]

	for _, r := range winners {
		if rank_list(r, orderedList) < rank_list(best, orderedList) {
			best = r
		}
	}

	return best
}

func rank_list(ranking []Alternative, ordereList [][]Alternative) int {
	i := 0
	for index, val := range ordereList {
		if equal_prefs(val, ranking) {
			i = index
			break
		}
	}

	return i
}
