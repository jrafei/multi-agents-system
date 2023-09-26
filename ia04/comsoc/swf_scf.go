package comsoc

func SWFFactory(swf func(p Profile) (Count, error), tieb func([]Alternative) (Alternative,error)) func(Profile) ([]Alternative, error) {
	f_swf := func(p Profile) ([]Alternative, error) {
		// Construction de la fonction avec application tiebreak
		count, err := swf(p)
		// Récupération du décompte
		if err != nil {
			return nil, err
		}
		alt, err := tieb(maxCount(count))
		// Application du tiebreak sur les objets ayant le plus de voix
		if err != nil {
			return nil, err
		}
		winning_alt := make([]Alternative,1)
		winning_alt[0] = alt
		// On construit une slice avec la solution
		return winning_alt, nil
	}
	return f_swf
}

func SCFFactory(scf func(p Profile) ([]Alternative, error), tieb func([]Alternative) (Alternative, error)) func(Profile) ([]Alternative, error) {
	f_scf := func(p Profile) ([]Alternative, error) {
		// Construction de la fonction avec application tiebreak
		alts, err := scf(p)
		// Récupération du décompte
		if err != nil {
			return nil, err
		}
		alt, err := tieb(alts)
		// Application du tiebreak sur les objets ayant le plus de voix
		if err != nil {
			return nil, err
		}
		winning_alt := make([]Alternative,1)
		winning_alt[0] = alt
		// On construit une slice avec la solution
		return winning_alt, nil
	}
	return f_scf
}
