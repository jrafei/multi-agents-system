package comsoc

import (
	"errors"
	"fmt"
)

/*--------------- TYPES DE BASE ---------------*/

type Alternative int
type Profile [][]Alternative
type Count map[Alternative]int

/********* TYPES AJOUTES *****************/

// Structure définissant un tuple d'alternative ("first bat second")
type AltTuple struct {
	first  Alternative
	second Alternative
}

func (t *AltTuple) First() Alternative {
	return t.first
}
func (t *AltTuple) Second() Alternative {
	return t.second
}


/*---------- FONCTIONS UTILITAIRES ------------*/

// renvoie l'indice ou se trouve alt dans prefs
func rank(alt Alternative, prefs []Alternative) int {
	i := 0
	for index, val := range prefs {
		if val == alt {
			i = index
			break
		}
	}
	return i
}

// renvoie vrai ssi alt1 est préférée à alt2
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	return rank(alt1, prefs) <= rank(alt2, prefs)
}

// renvoie les meilleures alternatives pour un décomtpe donné
func maxCount(count Count) (bestAlts []Alternative) {
	// Récupération des clés de valeur max ( plusieurs clés possibles )
	bestAlts = make([]Alternative, 0)
	var max_pts int
	for _, alt := range count {
		max_pts = alt
		break
	}
	for k, v := range count {
		if v == max_pts {
			// On ajoute la clé si elle est égale à la valeur max
			bestAlts = append(bestAlts, k)
		} else if v > max_pts {
			// On reconstruit un tableau d'une clé si plus grand
			bestAlts = make([]Alternative, 1)
			bestAlts[0] = k
			max_pts = v
		}
	}
	return bestAlts
}

// vérifie les préférences d'un agent, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois
func CheckProfile(prefs []Alternative, alts []Alternative) error {
	if len(prefs) == 0 {
		return fmt.Errorf("the list of preference is empty")
	}

	if len(alts) < len(prefs){
		return fmt.Errorf("there are more alternatives (%d) in preference than required (%d)",len(prefs),len(alts))
	}else if len(alts) > len(prefs){
		return fmt.Errorf("there are less alternatives (%d) in preference than required (%d)",len(prefs),len(alts))
	}

	// Verification que chaque alternative de la liste 'alts' apparait une seule fois dans les préférences
	for _, alt1 := range alts {
		cpt := 0
		for _, alt2 := range prefs {
			if alt1 == alt2 {
				cpt++
			}
			if cpt > 1 {
				return fmt.Errorf("alternative %d appears more than once", alt1)
			}
		}
		if cpt == 0 {
			return fmt.Errorf("alternative %d does not appear", alt1)
		}

	}
	return nil
}


// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func CheckProfileAlternative(prefs Profile, alts []Alternative) error {
	for _, pref := range prefs {
		if CheckProfile(pref, alts) != nil {
			return errors.New("profil is not valid")
		}
	}
	return nil
}
