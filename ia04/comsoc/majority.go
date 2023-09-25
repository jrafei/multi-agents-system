package comsoc

func MajoritySWF(p Profile) (count Count, err error) {
	err = CheckProfile(p) // à voir si on utilise CheckProfileAlternative()
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
}

func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	var count Count
	count, err = MajoritySWF(p)
	if err != nil {
		return nil,err
	}
	return maxCount(count),err
}

