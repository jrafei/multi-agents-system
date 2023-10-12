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

// Compte le nombre de fois que les alternatives du profil battent les autre alternatives.
// La fonction retourne un dictionnaire dont les clés sont des tuples (a,b) (" a bat b "), et les valeurs le nombre de fois que cela arrive.
func CountIsPref(p Profile) map[AltTuple]int {
	win := make(map[AltTuple]int) // enregistre le nombre de fois où a bat b
	for _, pref := range p {
		for index, alt := range pref {
			if alt == pref[len(pref)-1] {
				// On stop si on arrive à la dernière valeur (inutile de l'étudier car elle est battue par tout le monde)
				break
			} else {
				for _, alt2 := range pref[index+1:] {
					tuple := AltTuple{alt, alt2}
					_, exist := win[tuple]
					if exist {
						win[tuple]++
					} else {
						win[tuple] = 1
					}
				}
			}
		}
	}
	return win
}


// renvoie les meilleures alternatives pour un décomtpe donné
func minCount(count Count) (worstAlts []Alternative) {
	// Récupération des clés de valeur max ( plusieurs clés possibles )
	worstAlts = make([]Alternative, 0)
	var min_pts int
	for _, alt := range count {
		min_pts = alt
		break
	}
	for k, v := range count {
		if v == min_pts {
			// On ajoute la clé si elle est égale à la valeur max
			worstAlts = append(worstAlts, k)
		} else if v < min_pts {
			// On reconstruit un tableau d'une clé si plus grand
			worstAlts = make([]Alternative, 1)
			worstAlts[0] = k
			min_pts = v
		}
	}
	return worstAlts
}


// Elimination d'un élément, à partir de son index, dans une slice
func Remove(s []Alternative, index int) []Alternative {
	ret := make([]Alternative, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
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
// ********************************************************* A modifier ************************************
func CheckProfile(prefs []Alternative, alts []Alternative) error {
	if len(prefs) == 0 {
		return errors.New("the list of preference is empty")
	}

	// Verification que chaque alternative de la liste 'alts' apparait une seule fois dans les préférences
	for _, alt1 := range alts {
		cpt := 0
		for _, alt2 := range prefs {
			if alt1 == alt2 {
				cpt++
			}
			if cpt > 1 {
				return errors.New(fmt.Sprintf("alternative %d appears more than once", alt1))
			}
		}
		if cpt == 0 {
			return errors.New(fmt.Sprintf("alternative %d does not appear", alt1))
		}

	}
	return nil
}

// Vérifie que chaque alternative n'apparaît qu'une seule fois par préférence
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

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func CheckProfileAlternative(prefs Profile, alts []Alternative) error {
	for _, pref := range prefs {
		if CheckProfile(pref, alts) != nil {
			return errors.New("profil is not valid")
		}
	}
	return nil
}
