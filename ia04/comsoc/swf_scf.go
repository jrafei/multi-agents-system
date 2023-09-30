package comsoc

// Elimination d'un élément, à partir de son index, dans une slice
func remove(s []Alternative, index int) []Alternative {
	ret := make([]Alternative, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func SWFFactory(swf func(p Profile) (Count, error), tieb func([]Alternative) (Alternative, error)) func(Profile) ([]Alternative, error) {
	f_swf := func(p Profile) ([]Alternative, error) {
		// Construction de la fonction avec application tiebreak
		count, err := swf(p)
		// Récupération du décompte
		if err != nil {
			return nil, err
		}
		alts := maxCount(count)
		// Récupération des meilleurs alternatives
		sorted_alts := make([]Alternative, 0)
		for len(alts) > 0 {
			// Tri des alternatives en fonction du tiebreak
			alt, err := tieb(alts)
			if err != nil {
				return nil, err
			}
			index := rank(alt, alts)
			alts = remove(alts, index)
			sorted_alts = append(sorted_alts, alt)
		}
		return sorted_alts, nil
	}
	return f_swf
}

func SCFFactory(scf func(p Profile) ([]Alternative, error), tieb func([]Alternative) (Alternative, error)) func(Profile) (Alternative, error) {
	f_scf := func(p Profile) (Alternative, error) {
		// Construction de la fonction avec application tiebreak
		alts, err := scf(p)
		// Récupération du décompte
		if err != nil {
			return 0, err
		}
		alt, err := tieb(alts)
		// Application du tiebreak sur les objets ayant le plus de voix
		if err != nil {
			return 0, err
		}
		return alt, nil
	}
	return f_scf
}
