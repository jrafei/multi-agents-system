package comsoc

import (
	"errors"
)

/*--------------- TYPES DE BASE ---------------*/

type Alternative int
type Profile [][]Alternative
type Count map[Alternative]int

/********* TYPES AJOUTES *****************/

type AltTuple struct {
	first Alternative
	second Alternative
}

func (t *AltTuple) First() Alternative{
	return t.first
}
func (t *AltTuple) Second() Alternative{
	return t.second
}

/*---------- FONCTIONS UTILITAIRES ------------*/

// renvoie l'indice ou se trouve alt dans prefs
func rank(alt Alternative, prefs []Alternative) int {
	i := 0
	for index, val := range prefs {
		if val == alt {
			i = index
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
	max_pts := 0
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

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois par préférences
// à vérifier
func CheckProfile(prefs Profile) error {
	length := 0
	if len(prefs) > 0 {
		length = len(prefs[0])
	} else {
		return errors.New("no preference in profile")
	}
	for _, pref := range prefs {
		// Vérification que chaque préférence à le même nombre d'alternatives
		if len(pref) != length {
			return errors.New("not the same number of alternatives between preferences")
		}
		if checkAlternative(pref) != nil {
			return errors.New("alternative appears more than once in a preference")
		}
	}
	return nil
}

// Vérifie que chaque alternative n'apparaît qu'une seule fois par préférence
// à vérifier
func checkAlternative(pref []Alternative) error {
	check := make(map[Alternative]int) // nombre d'occurence des alternatives dans la préférence
	for _, v := range pref {
		_, present := check[v]
		if present {
			return errors.New("alternative appears more than once")
		} else {
			check[v] = 1
		}
	}
	return nil
}

func CheckProfileAlternative(prefs Profile, alts []Alternative) error {
	length := 0
	if len(prefs) > 0 {
		length = len(prefs[0])
	} else {
		return errors.New("no preference in profile")
	}
	for _, pref := range prefs {
		// Vérification que chaque préférence à le même nombre d'alternatives
		if len(pref) != length {
			return errors.New("not the same number of alternatives between preferences")
		}
		// Verification que chaque alternative n'apparait pas plusieur fois
		test := checkAlternative(pref)
		if test != nil {
			return errors.New("alternantive appears more than one")
		}

		//Verification que toutes les alternatives apparaissent dans la préference
		for _, alt1 := range alts {
			present := false
			for _, alt2 := range pref {
				if alt1 == alt2 {
					present = true
					break
				}
			}
			if !present {
				return errors.New("une alternative n'apparait pas dans la préférence")
			}
		}
	}
	return nil
}
