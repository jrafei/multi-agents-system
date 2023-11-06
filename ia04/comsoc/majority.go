package comsoc

import "errors"

//renvoie à partir d'un profile , le nombre de vote en faveur de chaque alternative
/*
======================================

	  @brief :
	  'Fonction de calcul du classement (SWF) de la méthode de vote par Majorité.'
	  @params :
		- 'p' : profile sur lequel appliquer la méthode
	  @returned :
	    -  'count' : le décompte des points
		- 'err' : erreur (nil si aucune erreur)

======================================
*/
func MajoritySWF(p Profile) (count Count, err error) {
	if len(p) == 0 {
		return nil, errors.New("profil is empty")
	}
	err = CheckProfileAlternative(p, p[0]) // à voir si on utilise CheckProfileAlternative()
	if err != nil {
		return nil, err
	}

	count = make(map[Alternative]int)
	// Initialisation des décomptes à 0
	for _, alt := range p[0] {
		count[alt] = 0
	}

	// Comptage
	for _, pref := range p {
		_, exist := count[pref[0]]
		if exist {
			count[pref[0]]++
		} else {
			count[pref[0]] = 1
		}
	}

	return count, nil
}

/*
======================================

	  @brief :
	  'Fonction de calcul du gagnant (SCF) de la méthode de vote par Majorité.'
	  @params :
		- 'p' : profile sur lequel appliquer la méthode
	  @returned :
	    -  'bestAlt' : le gagnant (vide si erreur)
		- 'err' : erreur (nil si aucune erreur)

======================================
*/
func MajoritySCF(p Profile) (bestAlt []Alternative, err error) {
	var count Count
	count, err = MajoritySWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), err
}
