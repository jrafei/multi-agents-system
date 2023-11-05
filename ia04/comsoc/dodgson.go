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
			list_new_prefs := flip_pref(pref, n_flips, nil) //renvoie une liste de préférences possibles après 'n_flips' de 'pref'

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
		return []Alternative{meilleurElement(liste_gagnant, orderedAlts)}, nil
	} else {
		return liste_gagnant, nil
	}
}

/*
renvoie la liste des preferences possibles après n inversion
avec n_flips >= 1
*/
func flip_pref(pref []Alternative, n_flips int, pere []Alternative) [][]Alternative {

	if n_flips == 1 {
		return one_flip(pref, pere)
	} else {
		res := one_flip(pref, pere)
		//fmt.Println("pour n_flips , ", n_flips)
		//fmt.Println("res = ", res)
		pref_possible := make([][]Alternative, 0)
		for _, y := range res {
			//fmt.Println("y = ", y)
			z := flip_pref(y, n_flips-1, pref)
			//fmt.Println("z = ", z)
			pref_possible = append(pref_possible, z...)
		}
		return pref_possible
	}

}

/*
renvoie la liste des preferences possibles après une inversion
*/
func one_flip(pref []Alternative, pere []Alternative) [][]Alternative {

	list_pref := make([][]Alternative, 0)
	
	for i := 0; i < len(pref)-1; i++ {
		copy_pref := make([]Alternative, len(pref))
		copy(copy_pref, pref)
		copy_pref[i] = pref[i+1]
		copy_pref[i+1] = pref[i]

		if len(pere) == 0 || !equal_prefs(pere, copy_pref) {
			list_pref = append(list_pref, copy_pref)
		}

	}

	return list_pref
}
