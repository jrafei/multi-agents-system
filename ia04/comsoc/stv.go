package comsoc

import (
	"errors"
	"fmt"
)

func STV_SWF(p Profile,orderedAlts []Alternative) (count Count, err error) {
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

	counter := len(p[0])-1 // récupération du nb de préférences (car au pire on fait nb_de_prefs tests de majorité)
	for counter > 0 {
		fmt.Println(p)
		maj_count, err := MajoritySWF(p) // Réalisation du test de majorité
		if err != nil {
			return nil, err
		}
		if absoluteMajority(p, maj_count) {
			// Si la majorité absolue est atteinte, on peut directement retourner les valeurs
			count[maxCount(maj_count)[0]]++
			fmt.Println("ok")
			fmt.Println(count)
			return count, nil
		} else {

			worstAlts := minCount(maj_count)								// Récupération des pires alternatives
			reversedOrderedAlts := inversion(orderedAlts) 					// Inversion du tiebreak, pour réutilisation de la factory
			worst,err := TieBreakFactory(reversedOrderedAlts)(worstAlts) 		// Application du tiebreak
			if err != nil{
				return nil, err
				}
			p = removeProfile(p,worst)										// On supprime la pire alternative dans le profile

			fmt.Println(maj_count)
			fmt.Println(worst)

			// On met à jour count, en recupérant les alternatives utilisées dans le test
			for alt,_ := range maj_count {
				if(alt!=worst){
					count[alt]++
				}
			}

			
			fmt.Println(count)
			}
			counter--	// On indique qu'on a supprimé une alternative
			
		}
	
	return count, nil

}

func absoluteMajority(p Profile, count Count) bool {
	// Vérification si la majorité absolue est atteinte
	maj_abs := (len(p) / 2) + 1
	for _, votes := range count {
		if votes > maj_abs {
			return true
		}
	}
	return false
}

func inversion(ordered []Alternative)(inverted []Alternative){
	length := len(ordered) 
	inverted = make([]Alternative,length)
	for i:=length-1;i>=0;i--{
		inverted[length-i-1] = ordered[i]
	}
	return inverted
}

func removeProfile(p Profile, alt Alternative) (new_p Profile){
	for i,pref := range p {
		p[i] = Remove(pref,rank(alt,pref))
	}
	return p
}

// func STV_SCF(p Profile) (bestAlts []Alternative, err error) {}
