package comsoc

import (
	"errors"
	"fmt"
)

func YoungCondorcet(p Profile, orderedAlts []Alternative) (bestAlt []Alternative, err error) {
	if len(p) == 0 {
		return nil, errors.New("profil is empty")
	}

	if len(orderedAlts) == 0 {
		return nil, errors.New("tiebreak is empty")
	}

	err = CheckProfileAlternative(p, p[0])
	if err != nil {
		return nil, err
	}

	tmp_bestAlt := make([]Alternative, 0)
	nb_removed_pref := 0

	for nb_removed_pref <= len(p) {
		for _, combination := range combinationsKamongN(len(p), nb_removed_pref) {
			// Détermine toutes les combinaisons possibles de préférences à supprimer
			tmp_profile := p
			for i, pref_index := range combination {
				// Suppression de la préférence
				tmp_profile = RemovePref(tmp_profile, pref_index-1-i)
			}
			winner, err := CondorcetWinner(tmp_profile)
			if err != nil {
				return nil, err
			} else if len(winner) == 1 {
				if winner[0] == orderedAlts[0] {
					// Si on a trouvé un gagant de Condorcet, et que c'est le préféré. Il a automatiquement gagné
					return winner, nil
				} else {
					tmp_bestAlt = append(bestAlt, winner[0])
				}
			}
			fmt.Println(tmp_profile)
		}
		// Si bestAlt est rempli on retourne la meilleur alternative
		if len(tmp_bestAlt) == 1 {
			// Si qu'un seul gagnant pour le même nombre de préférences supprimées, il n'y a pas d'ambuiguité
			return tmp_bestAlt, nil
		} else if len(tmp_bestAlt) > 1 {
			return []Alternative{meilleurElement(tmp_bestAlt, orderedAlts)},nil
		} else {
			// Sinon, on teste en supprimant plus de préférences
			nb_removed_pref++
		}
		
	}
	// Cas normalement impossible à atteindre (lorsqu'il reste qu'un seul vote, il y a forcément un gagnant de Condorcet)
	return nil, errors.New("Aucun gagnant Young-Condorcet")

}



// Retourne le meilleur élément
func meilleurElement(elements []Alternative, classement []Alternative) Alternative {
	best := elements[0]
	for _,alt := range elements{
		if rank(alt,classement) < rank(best,classement) {
			best = alt
		}
	}
	return best
}

// Elimination d'un élément, à partir de son index, dans une slice
func RemovePref(s [][]Alternative, index int) [][]Alternative {
	ret := make([][]Alternative, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func combinationsKamongN(n, k int) [][]int {
	if k > n {
		return [][]int{}
	}

	var result [][]int

	// Fonction récursive pour générer les combinaisons
	var generateCombinations func(start, count int, combination []int)

	generateCombinations = func(start, count int, combination []int) {
		if count == 0 {
			// Une combinaison est prête, ajoutez-la au résultat
			temp := make([]int, len(combination))
			copy(temp, combination)
			result = append(result, temp)
			return
		}

		for i := start; i <= n-count+1; i++ {
			// Ajoutez l'élément actuel à la combinaison
			combination[len(combination)-count] = i
			// Générer les combinaisons restantes
			generateCombinations(i+1, count-1, combination)
		}
	}

	// Initialisez la slice de combinaison avec la taille k
	combination := make([]int, k)

	// Commencez à générer les combinaisons
	generateCombinations(1, k, combination)

	return result
}
