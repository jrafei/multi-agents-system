package comsoc


import (
	"errors"
)

func TieBreakFactory(orderedAlts []Alternative) (func ([]Alternative) (Alternative, error)){
	tiebrake := func(alts []Alternative) (a Alternative, err error) {
		if len(alts)==0 {
			return -1,errors.New("no alternative to order")
		}

		if len(orderedAlts)==0 {
			return alts[0],errors.New("unable to order alternatives")
		}

		if len(alts)==1{
			return alts[0],nil
		}

		winning_alt:= alts[0]
		for pos := range alts{
			if rank(alts[pos],orderedAlts) < rank(winning_alt,orderedAlts){
				winning_alt = alts[pos]
			}
		}
		return winning_alt,nil
	}
	return tiebrake
}

