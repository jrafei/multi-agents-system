package comsoc


func BordaSWF(p Profile) (count Count, err error){
	err = CheckProfile(p) // à voir si on utilise CheckProfileAlternative()
	if err != nil {
		return nil, err
	}
	count = make(map[Alternative]int)
	for _, pref := range p {
		for index,key := range pref{
			_, exist := count[key]
			if exist {
				count[key] = count[key]+(len(pref)-index-1)
			} else {
				count[key] = len(pref)-index-1
			}
		}
	}
	return count,nil
}
func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	var count Count
	count, err = BordaSWF(p)
	if err != nil {
		return nil,err
	}
	return maxCount(count),err
}
