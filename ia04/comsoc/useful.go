package comsoc

/*
======================================
	*
	*
	*  Ce fichier déclare et définie les fonctions annexes utilisées pour nos algorithmes de calcul des classements et gagnants.
	*
	*
======================================
*/

/*
======================================

	  @brief :
	  'Meilleure alternative d une liste d alternatives, selon un classement.'
	  @params :
		- 'elements' : liste des éléments à etudier
		- 'classement' : classement des alternatives
	  @returned :
	    -  La meilleure alternative selon le classement.

======================================
*/
func meilleurElement(elements []Alternative, classement []Alternative) Alternative {
	best := elements[0]
	for _, alt := range elements {
		if rank(alt, classement) < rank(best, classement) {
			best = alt
		}
	}
	return best
}

/*
======================================

	  @brief :
	  'Elimination d'une préférence d un profil, à partir de son index.'
	  @params :
		- 'elements' : liste des éléments à etudier
		- 'index' : la position de la préférence dans le profil
	  @returned :
	    -  La meilleure alternative selon le classement.

======================================
*/
func removePref(s [][]Alternative, index int) [][]Alternative {
	ret := make([][]Alternative, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

/*
======================================

	  @brief :
	  'Ensemble des combinaisons de k parmi n.'
	  @params :
		- 'n' : nombre total d éléments
		- 'k' : nombre souhaité d éléments
	  @returned :
	    -  La meilleure alternative selon le classement.

======================================
*/
func combinationsKamongN(n, k int) [][]int {
	if k > n {
		return [][]int{}
	}

	var result [][]int

	// Fonction récursive pour générer les combinaisons
	var generateCombinations func(start, count int, combination []int)

	generateCombinations = func(start, count int, combination []int) {
		if count == 0 {
			temp := make([]int, len(combination))
			copy(temp, combination)
			result = append(result, temp)
			return
		}

		for i := start; i <= n-count+1; i++ {
			combination[len(combination)-count] = i
			generateCombinations(i+1, count-1, combination)
		}
	}

	combination := make([]int, k)

	generateCombinations(1, k, combination)

	return result
}

/*
======================================

	  @brief :
	  'Compte le nombre de fois que les alternatives du profil battent les autres alternatives.'
	  @params :
		- 'p' : profil à étudier
	  @returned :
	    -  La fonction retourne un dictionnaire dont les clés sont des tuples (a,b) (" a bat b "), et les valeurs le nombre de fois que cela arrive.

======================================
*/
func countIsPref(p Profile) map[AltTuple]int {
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

//
/*
======================================

	  @brief :
	  'Renvoie les pires alternatives pour un décompte donné (celles avec le moins de points).'
	  @params :
		- 'count' : le décompte
	  @returned :
		- 'worstAlts' : les pires alternatives
======================================
*/
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
			// On ajoute la clé si elle est égale à la valeur min
			worstAlts = append(worstAlts, k)
		} else if v < min_pts {
			// On reconstruit un tableau d'une clé si plus petit
			worstAlts = make([]Alternative, 1)
			worstAlts[0] = k
			min_pts = v
		}
	}
	return worstAlts
}

/*
======================================

	  @brief :
	  'Elimination d une alternative, à partir de son index dans une préférence.'
	  @params :
		- 's' : la préférence
		- 'index' : la position de l'alternative à supprimer
	  @returned :
	    -  La meilleure alternative selon le classement.

======================================
*/
func removeAlt(s []Alternative, index int) []Alternative {
	ret := make([]Alternative, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

/*
======================================

	  @brief :
	  'Compare deux préférences.'
	  @params :
		- 'pref1' : la préférence 1
		- 'pref2' : la préférence 2
	  @returned :
	    -  booléen vrai si égalité

======================================
*/
func equal_prefs(pref1 []Alternative, pref2 []Alternative) bool {

	for k, alt1 := range pref1 {
		if alt1 != pref2[k] {
			return false
		}

	}
	return true
}

/*
======================================

	  @brief :
	  'Echange deux éléments d une liste d alternatives'
	  @params :
	  	- 'numbers' : liste d alternatives
		- 'i' : indice de l élément 1
		- 'j' : indice de l élément 2

======================================
*/
func swap(numbers []Alternative, i, j int) {
	if i >= len(numbers) || j >= len(numbers) || i < 0 || j < 0 || i == j {
		return
	}
	numbers[i], numbers[j] = numbers[j], numbers[i]
}

/*
Fonction récursive pour générer les permutations
*/

/*
======================================

	  @brief :
	  'Génération des permutations possibles d une liste d'alternatives.'
	  @params :
	  	- 'numbers' : liste d alternatives
		- 'start' : indice de départ
		- 'result' : pointeur de la liste des permutations possibles

======================================
*/
func permute(numbers []Alternative, start int, result *[][]Alternative) {
	if start == len(numbers)-1 {
		// Fait une copie de la permutation courante pour ne pas modifier le résultat
		perm := make([]Alternative, len(numbers))
		copy(perm, numbers)
		*result = append(*result, perm)
		return
	}

	for i := start; i < len(numbers); i++ {
		// Échange le début avec l'élément courant
		swap(numbers, start, i)
		// Appel récursif pour les éléments restants
		permute(numbers, start+1, result)
		// Restaure l'ordre initial pour le prochain itération
		swap(numbers, start, i)
	}
}

/*
======================================

	  @brief :
	  'Vérification si la majorité absolue est atteinte pour le décompte d un profil donnée (ATTENTION, cela suppose que : un point = une vote !).'
	  @params :
	    - 'numbers' : liste d alternatives
		- 'start' : indice de départ
		- 'result' : pointeur de la liste des permutations possibles
	  @returned :
	  	- booléen vrai si majorité atteinte.

======================================
*/
func absoluteMajority(p Profile, count Count) bool {
	// Vérification si la majorité absolue est atteinte
	maj_abs := (len(p) / 2) + 1
	for _, votes := range count {
		if votes >= maj_abs {
			return true
		}
	}
	return false
}

/*
======================================

	  @brief :
	  'Inverse un tableau d alternatives.'
	  @params :
	  	- 'ordered' : tableau à inverser
	  @returned :
	  	- 'inverted' : tableau inversé

======================================
*/
func inversion(ordered []Alternative) (inverted []Alternative) {
	length := len(ordered)
	inverted = make([]Alternative, length)
	for i := length - 1; i >= 0; i-- {
		inverted[length-i-1] = ordered[i]
	}
	return inverted
}

/*
======================================

	  @brief :
	  'Inverse un tableau d alternatives.'
	  @params :
        - 'p' : profil à étudier
		- 'alt' : alternative à supprimer
	  @returned :
	  	- 'new_p' : nouveau profil

======================================
*/
func removeAltProfile(p Profile, alt Alternative) (new_p Profile) {
	for i, pref := range p {
		p[i] = removeAlt(pref, rank(alt, pref))
	}
	return p
}

/*
======================================

	  @brief :
	  'Calcul du score de classement de Kemeny-Young'
	  @params :
	    - 'ranking' : classement
		- 'alt' :
	  @returned :
	  	- score du classement

======================================
*/
func calculateScoreKemenyYoung(ranking []Alternative, battle map[AltTuple]int) int {
	res := 0
	for x, _ := range ranking {
		for y := x + 1; y < len(ranking); y++ {
			res += battle[AltTuple{Alternative(ranking[x]), Alternative(ranking[y])}]
		}
	}
	return res

}

/*
======================================

	  @brief :
	  'Renvoie la liste des preferences possibles après n inversions avec n_flips >= 1'
	  @params :
	    - 'pref' : préférence de départ
		- 'n_flips' : nombre de flips à effectuer
		- 'pere' : préférence parente (récursivité)
	  @returned :
	  	- liste des préférences avec nouveaux flips

======================================
*/
func flip_pref(pref []Alternative, n_flips int, pere []Alternative) [][]Alternative {

	if n_flips == 1 {
		return one_flip(pref, pere)
	} else {
		res := one_flip(pref, pere)
		pref_possible := make([][]Alternative, 0)
		for _, y := range res {
			z := flip_pref(y, n_flips-1, pref)
			pref_possible = append(pref_possible, z...)
		}
		return pref_possible
	}

}

/*
======================================

	  @brief :
	  'Renvoie la liste des preferences possibles après une inversion'
	  @params :
        - 'pref' : préférence de départ
		- 'pere' : préférence parente (récursivité)
	  @returned :
	  	- liste des préférences avec nouveaux flips

======================================
*/
func one_flip(pref []Alternative, pere []Alternative) [][]Alternative {

	list_pref := make([][]Alternative, 0)

	for i := 0; i < len(pref)-1; i++ {
		copy_pref := make([]Alternative, len(pref))
		copy(copy_pref, pref)
		copy_pref[i] = pref[i+1]
		copy_pref[i+1] = pref[i]

		if len(pere) == 0 || !equal_prefs(pere, copy_pref) {
			list_pref = append(list_pref, copy_pref)
		}

	}

	return list_pref
}
