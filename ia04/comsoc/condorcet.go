package comsoc

import (
	"errors"
)

/*
======================================

	  @brief :
	  'Calcul du gagnant de Condorcet'
	  @params :
		- 'p' : profile sur lequel appliquer la méthode
	  @returned :
	    -  'bestAlts' : gagnant de la méthode (vide si aucun gagnant, de taille 1 sinon)
		- 'err' : erreur (nil si aucune erreur)

======================================
*/
func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
 
	if len(p) == 0 {
		return nil, errors.New("no preference in profile ")
	}

	N := len(p[0])
	if N == 0 {
		return nil, errors.New("preference is empty")
	}

	err = CheckProfileAlternative(p, p[0])
	if err != nil {
		return nil, err
	}

	list := p[0]
	for _, alt1 := range list {
		battre := true
		for _, alt2 := range list {
			if alt1 != alt2 {
				count1 := 0
				count2 := 0
				for _, prefs := range p {
					if isPref(alt1, alt2, prefs) {
						count1++
					} else {
						count2++
					}
				}
				if count1 < count2 {
					battre = false
					break
				}
			}
		}
		if battre {
			bestAlts = append(bestAlts, alt1)
		}
	}
	if len(bestAlts) > 1 {
		// AUCUN gagnant de condorcet en cas de tie, on retourne une slice vide
		return make([]Alternative, 0), nil
	}

	return bestAlts, nil
}
