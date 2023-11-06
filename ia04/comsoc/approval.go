package comsoc

import (
	"errors"
)

/*
======================================

	  @brief :
	  'Fonction de calcul du classement (SWF) de la méthode de vote par Approbation.'
	  @params :
		- 'p' : profile sur lequel appliquer la méthode
		- 'thresholds' : tableau des seuils de vote des élécteurs
	  @returned :
	    -  'count' : le décompte des points 
		- 'err' : erreur (nil si aucune erreur)

======================================
*/
func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	if len(p) == 0 {
		return nil, errors.New("profil is empty")
	}
	err = CheckProfileAlternative(p, p[0])
	if err != nil {
		return nil, err
	}

	if len(thresholds) != len(p) {
		return nil, errors.New("the number of thresholds doesn't match the profile")
	}

	count = make(map[Alternative]int)
	// Initialisation des décomptes à 0
	for _, alt := range p[0] {
		count[alt] = 0
	}

	for index_profile, pref := range p {
		if thresholds[index_profile] > len(pref) {
			return nil, errors.New("the thresholds exceeds the preference length")
		}
		for _, key := range pref[:thresholds[index_profile]] {
			// On itère uniquement entre l'indice 0 et le seuil associé (indice exclu)
			_, exist := count[key]
			if exist {
				count[key]++
			} else {
				count[key] = 1
			}
		}
	}
	return count, nil
}

/*
======================================

	  @brief :
	  'Fonction de calcul du gagant (SCF) de la méthode de vote par Approbation.'
	  @params :
		- 'p' : profile sur lequel appliquer la méthode
		- 'thresholds' : tableau des seuils de vote des élécteurs
	  @returned :
	    -  'bestAlt' : le gagnant (vide si erreur)
		- 'err' : erreur (nil si aucune erreur)

======================================
*/
func ApprovalSCF(p Profile, thresholds []int) (bestAlt []Alternative, err error) {
	var count Count
	count, err = ApprovalSWF(p, thresholds)
	if err != nil {
		return nil, err
	}
	return maxCount(count), err
}
