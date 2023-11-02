package comsoc

import(
	"errors"
)


func KramerSimpson_SWF(p Profile) (count Count, err error){
	if len(p) == 0 {
		return nil, errors.New("profil is empty")
	}

	err = CheckProfileAlternative(p, p[0])
	if err != nil {
		return nil, err
	}
	count = make(Count,0)
	nbAts := len(p[0])
	i := 1
	for i<=nbAts{
		duels := CountIsPref(p)
		min_val_duel := len(p)
		for tuple,value := range duels{
			if tuple.first == Alternative(i) && value < min_val_duel {
				min_val_duel = value
			}
		}
		count[Alternative(i)] = min_val_duel
		i++;
	}
	return
}

func KramerSimpson_SCF(p Profile) (bestAlts []Alternative, err error) {
	var count Count
	count, err = KramerSimpson_SWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), err
}