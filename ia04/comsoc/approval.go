package comsoc

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	err = CheckProfile(p) // à voir quelle(s) vérifications on doit faire
	if err != nil {
		return nil, err
	}
	count = make(map[Alternative]int)
	for index_profile, pref := range p {
		for _, key := range pref[:thresholds[index_profile]] {
			// On itère uniquement entre l'indice 0 et le seuil associé (indice exclu)
			_, exist := count[key]
			if exist {
				count[key]++
			} else {
				count[key] = 1
			}
		}
	}
	return count, nil
}
func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	var count Count
	count, err = ApprovalSWF(p, thresholds)
	if err != nil {
		return nil, err
	}
	return maxCount(count), err
}
