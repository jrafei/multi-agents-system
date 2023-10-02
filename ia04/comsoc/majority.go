package comsoc

import "errors"

//renvoie à partir d'un profile , le nombre de vote en faveur de chaque alternative
func MajoritySWF(p Profile) (count Count, err error) {
	if len(p) == 0 {
		return nil, errors.New("profil is empty")
	}
	err = CheckProfileAlternative(p, p[0]) // à voir si on utilise CheckProfileAlternative()
	if err != nil {
		return nil, err
	}

	count = make(map[Alternative]int)
	for _, pref := range p {
		_, exist := count[pref[0]]
		if exist {
			count[pref[0]]++
		} else {
			count[pref[0]] = 1
		}
	}

	return count, nil
	// peut-être retourner aussi les alts non comptées !
}

//renvoie à partir d'un profile, les alternantives qui ont un nombre de vote maximal

func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	var count Count
	count, err = MajoritySWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), err
}
