package comsoc

import (
	"errors"
)

func STV_SWF(p Profile) (count Count, err error) {
	if len(p) == 0 {
		return nil, errors.New("profil is empty")
	}
	err = CheckProfileAlternative(p, p[0])
	if err != nil {
		return nil, err
	}

	counter := len(p[0]) - 1
	count = make(Count)
	for _, alt := range p[0] {
		count[alt] = 0 // Initialisation du comptage à 0 pour chaque alternative
	}

	for counter > 0 {
		maj_count, err := MajoritySWF(p) // Réalisation du test de majorité
		if err != nil {
			return nil, err
		}
		if absoluteMajority(p, maj_count) {
			// Si la majorité absolue est atteinte, on peut directement retourner les valeurs
			count[maxCount(maj_count)[0]] = 1
			return count, nil
		} else {
			// IL FAUT VOIR SI ON SUPPRIME TOUTES LES ALTERNATIVES QUI N'ONT PAS ETE VOTEES, SINON LAQUELLE
		}

		counter--
	}
	return count, nil

}

func absoluteMajority(p Profile, count Count) bool {
	// Vérification si la majorité absolue est atteinte
	maj_abs := (len(p) / 2) + 1
	for _, votes := range count {
		if votes > maj_abs {
			return true
		}
	}
	return false
}

// func STV_SCF(p Profile) (bestAlts []Alternative, err error) {}
