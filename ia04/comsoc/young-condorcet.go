package comsoc

import (
	"errors"
	"fmt"
)

/*
======================================

	  @brief :
	  'Calcul du gagnant de Young-Condorcet'
	  @params :
		- 'p' : profile sur lequel appliquer la méthode
		- 'orderedAlts' : tiebreak pour le départage des alternatives
	  @returned :
	    -  'bestAlts' : gagnant de la méthode (vide si aucun gagnant, de taille 1 sinon)
		- 'err' : erreur (nil si aucune erreur)

======================================
*/
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
				tmp_profile = removePref(tmp_profile, pref_index-1-i)
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
			return []Alternative{meilleurElement(tmp_bestAlt, orderedAlts)}, nil
		} else {
			// Sinon, on teste en supprimant plus de préférences
			nb_removed_pref++
		}

	}
	// Cas normalement impossible à atteindre (lorsqu'il reste qu'un seul vote, il y a forcément un gagnant de Condorcet)
	return nil, errors.New("aucun gagnant Young-Condorcet")

}
