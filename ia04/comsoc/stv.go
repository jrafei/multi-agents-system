package comsoc

import (
	"errors"
)

/*
======================================

	  @brief :
	  'Fonction de calcul du classement (SWF) de la méthode de vote STV.'
	  @params :
		- 'p' : profile sur lequel appliquer la méthode
		- 'orderedAlts' : tiebreak pour le départage des alternatives
	  @returned :
	    -  'count' : le décompte des points 
		- 'err' : erreur (nil si aucune erreur)

======================================
*/
func STV_SWF(p Profile, orderedAlts []Alternative) (count Count, err error) {
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

	count = make(Count)
	for _, alt := range p[0] {
		count[alt] = 0 // Initialisation du comptage à 0 pour chaque alternative
	}

	counter := len(p[0]) - 1 // récupération du nb de préférences (car au pire on fait nb_de_prefs tests de majorité)
	for counter > 0 {

		maj_count, err := MajoritySWF(p) // Réalisation du test de majorité
		if err != nil {
			return nil, err
		}
		if AbsoluteMajority(p, maj_count) {
			// Si la majorité absolue est atteinte, on peut directement retourner les valeurs
			count[maxCount(maj_count)[0]]++
			return count, nil
		} else {

			worstAlts := minCount(maj_count)                              // Récupération des pires alternatives
			reversedOrderedAlts := Inversion(orderedAlts)                 // Inversion du tiebreak, pour réutilisation de la factory
			worst, err := TieBreakFactory(reversedOrderedAlts)(worstAlts) // Application du tiebreak
			if err != nil {
				return nil, err
			}
			p = RemoveAltProfile(p, worst) // On supprime la pire alternative dans le profile

			// On met à jour count, en recupérant les alternatives utilisées dans le test
			for alt, _ := range maj_count {
				if alt != worst {
					count[alt]++
				}
			}

		}
		counter-- // On indique qu'on a supprimé une alternative

	}

	return count, nil

}

/*
======================================

	  @brief :
	  'Fonction de calcul du gagnant (SCF) de la méthode de vote par Majorité.'
	  @params :
		- 'p' : profile sur lequel appliquer la méthode
		- 'orderedAlts' : tiebreak
	  @returned :
	    -  'bestAlt' : le gagnant (vide si erreur)
		- 'err' : erreur (nil si aucune erreur)

======================================
*/
func STV_SCF(p Profile, orderedAlts []Alternative) (bestAlt []Alternative, err error) {
	var count Count
	count, err = STV_SWF(p, orderedAlts)
	if err != nil {
		return nil, err
	}
	return maxCount(count), err
}
