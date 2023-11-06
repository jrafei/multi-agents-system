package comsoc

import "errors"

/*
 Reference du méthode : https://en.wikipedia.org/wiki/Kemeny%E2%80%93Young_method
*/

/*
======================================

	  @brief :
	  'Fonction de calcul du classement (SWF) de la méthode de vote de Kemeny-Young.'
	  @params :
		- 'p' : profile sur lequel appliquer la méthode
		- 'orderedAlts' : tiebreak pour le départage des alternatives
	  @returned :
	    -  'count' : le décompte des points 
		- 'err' : erreur (nil si aucune erreur)

======================================
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
	Permute(p[0], 0, &rankings)

	bestRank := [][]Alternative{rankings[0]}
	bestScore := CalculateScoreKemenyYoung(bestRank[0], battle)

	for i := 1; i < len(rankings); i++ {
		r := rankings[i]
		s := CalculateScoreKemenyYoung(r, battle)
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
		return []Alternative{MeilleurElement(bestWinners, orderedAlts)}, nil
	}

	//return []Alternative{bestRank[0][0]}, nil
	return bestRank[0], nil
}
