package comsoc

import "errors"

//Calcule le score de chaque alternative par la méthode de borda

func BordaSWF(p Profile) (count Count, err error) {
	if len(p) == 0 {
		return nil, errors.New("profil is empty")
	}
	err = CheckProfileAlternative(p, p[0]) // à voir si on utilise CheckProfileAlternative()
	if err != nil {
		return nil, err
	}
	count = make(map[Alternative]int)
	for _, pref := range p {
		for index, key := range pref {
			_, exist := count[key]
			if exist {
				count[key] = count[key] + (len(pref) - index - 1)
			} else {
				count[key] = len(pref) - index - 1
			}
		}
	}
	return count, nil
}

//renvoie les alternatives qui ont un score Borda maximal
func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	var count Count
	count, err = BordaSWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), err
}
