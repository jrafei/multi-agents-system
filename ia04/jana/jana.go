package jana

import (
	"errors"
)

/*--------------- TYPES DE BASE ---------------*/

type Alternative int
type Profile [][]Alternative
type Count map[Alternative]int

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
func checkProfile(prefs Profile) error {
	length := 0
	if len(prefs) > 0 {
		length = len(prefs[0])
	} else {
		return errors.New("No preference in profile.\n")
	}
	for _, pref := range prefs {
		// Vérification que chaque préférence à le même nombre d'alternatives
		if len(pref) != length {
			return errors.New("Not the same number of alternatives between preferences.")
		}
		if checkAlternative(pref) != nil {
			return errors.New("Alternative appears more than once in a preference.")
		}
	}
	return nil
}

// Vérifie que chaque alternative n'apparaît qu'une seule fois par préférence
// à vérifier
func checkAlternative(pref []Alternative) error {
	check := make(map[Alternative]int) // nombre d'occurence des alternatives dans la préférence
	for _, v := range pref {
		if check[v] == 0 {
			check[v] = 1
		} else if check[v] > 0 {
			return errors.New("Alternative appears more than once.")
		}
	}
	return nil
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func checkProfileAlternative_v2(prefs Profile, alts []Alternative) error {
	length := 0
	if len(prefs) > 0 {
		length = len(prefs[0])
	} else {
		return errors.New("No preference in profile.\n")
	}
	for _, pref := range prefs {
		// Vérification que chaque préférence à le même nombre d'alternatives
		if len(pref) != length {
			return errors.New("Not the same number of alternatives between preferences.")
		}
		// Verification que chaque alternative n'apparait pas plusieur fois
		test := checkAlternative(pref)
		if test != nil {
			return errors.New("Alternantive appears more than one")
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
			if present == false {
				return errors.New("Une alternative n'apparait pas dans la préférence")
			}
		}
	}
	return nil
}
