package comsoc

import (
	"errors"
)

/*
======================================

	  @brief :
	  'Fonction de calcul du classement (SWF) de la méthode de vote de Kramer-Simpson.'
	  @params :
		- 'p' : profile sur lequel appliquer la méthode
	  @returned :
	    -  'count' : le décompte des points
		- 'err' : erreur (nil si aucune erreur)

======================================
*/
func KramerSimpson_SWF(p Profile) (count Count, err error) {
	if len(p) == 0 {
		return nil, errors.New("profil is empty")
	}

	err = CheckProfileAlternative(p, p[0])
	if err != nil {
		return nil, err
	}
	count = make(Count, 0)
	nbAts := len(p[0])
	i := 1
	for i <= nbAts {
		duels := CountIsPref(p)
		min_val_duel := len(p)
		for tuple, value := range duels {
			if tuple.first == Alternative(i) && value < min_val_duel {
				min_val_duel = value
			}
		}
		count[Alternative(i)] = min_val_duel
		i++
	}
	return count, nil
}

/*
======================================

	  @brief :
	  'Fonction de calcul du gagnant (SCF) de la méthode de vote de Kramer-Simpson.'
	  @params :
		- 'p' : profile sur lequel appliquer la méthode
	  @returned :
	    -  'bestAlt' : le gagnant (vide si erreur)
		- 'err' : erreur (nil si aucune erreur)

======================================
*/
func KramerSimpson_SCF(p Profile) (bestAlt []Alternative, err error) {
	var count Count
	count, err = KramerSimpson_SWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), err
}
