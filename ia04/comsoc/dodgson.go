package comsoc

import (
	"errors"
)

/*
	Le meilleur candidat est celui qui requiert le moins de
	changement dans les préférences des individus pour devenir un gagnant
	de condorcet. (nombre de “flips” : inversion de la préférence entre 2
	candidats dans les préférences d’un individu)
*/

func Dodgson(p Profile, orderedAlts []Alternative) ([]Alternative, error) {
	if len(p) == 0 {
		return nil, errors.New("profil is empty")
	}

	err := CheckProfileAlternative(p, p[0])
	if err != nil {
		return nil, err
	}

	n_flips := 1

	liste_gagnant := make([]Alternative, 0) // les gagnants qui ont le moins de changement dans les préférences

	for {

		for k1, pref := range p {
			// on copie le Profile pour chaque traitement d'une préference d'un individu
			copy_p := make([][]Alternative, len(p))
			copy(copy_p, p)
			list_new_prefs := Flip_pref(pref, n_flips, nil) //renvoie une liste de préférences possibles après 'n_flips' de 'pref'

			for _, pref_possible := range list_new_prefs {
				// on échange l'ancienne preference par la nouvelle
				copy_p[k1] = pref_possible
				// on applique Condorcet sur le nouveau profil
				bestAlts, _ := CondorcetWinner(copy_p)

				// si le gagnant existe on l'ajoute dans le map des gagnants
				if len(bestAlts) == 1 {
					liste_gagnant = append(liste_gagnant, bestAlts[0])
				}
			}
		}

		// si on trouve au moins un gagnant pour le nombre de flip donnée, on sort de boucle sinon on cherche un gagnant pour le nombre de flip suivant
		// ce code va s'executer normalement à un moment , puisqu'il existe au moins un nombre de flips pourlequel il y a un gagnant Condorcet
		if len(liste_gagnant) != 0 {
			break
		}
		// aucun gagnant -> cherchons pour le nombre de flips suivant
		n_flips++
		liste_gagnant = make([]Alternative, 0)
	}

	if len(liste_gagnant) != 1 { // application de tiebreak
		return []Alternative{MeilleurElement(liste_gagnant, orderedAlts)}, nil
	} else {
		return liste_gagnant, nil
	}
}

