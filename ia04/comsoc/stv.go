package comsoc

func STV_SWF(p Profile) (count Count, err error) {

	alts := make([]Alternative, 0)

	for _, val := range p[0] {
		alts = append(alts, val)
	}
	err = CheckProfileAlternative(p, alts)
	// vérification des préferences, si une incohérence survient par rapport à la liste d'alternatives, une erreur sera générée
	if err != nil {
		return nil, err
	}
	counter := len(alts) - 1
	count = make(Count)
	for _, alt := range alts {
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
	maj_abs := len(p) //2
	for _, votes := range count {
		if votes > maj_abs {
			return true
		}
	}
	return false
}

// func STV_SCF(p Profile) (bestAlts []Alternative, err error) {}
