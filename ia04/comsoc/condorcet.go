package comsoc
/*
import (
	"errors"
)


func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	if len(p) == 0 {
		return nil, errors.New("no preference in profile ")
	}

	N := len(p[0])
	if N == 0 {
		return nil, errors.New("preference is empty")
	}

	list := combin.Combinations(N, 2)

	count := make(map[Alternative]int)
	for _, comb := range list {
		a := Alternative(comb[0] + 1) //à revoir, je considère que les alternatives sont de 1 à N
		b := Alternative(comb[1] + 1)
		cpta := 0
		cptb := 0
		for _, pref := range p {
			if isPref(a, b, pref) {
				cpta++
			} else {
				cptb++
			}
		}
		if cpta > cptb { //incrémenter la valeur d'alt a dans le map count s'il existe
			_, exist := count[a]
			if exist {
				count[a]++
			} else {
				count[a] = 1
			}
		}

		if cpta < cptb { //incrémenter la valeur d'alt a dans le map count s'il existe
			_, exist := count[a]
			if exist {
				count[b]++
			} else {
				count[b] = 1
			}
		}
	}

	bestAlts = maxCount(count)
	if len(bestAlts) > 1 {
		return nil, errors.New("best alternative doesn't exist")
	}
	return bestAlts, nil
}
*/