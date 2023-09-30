package main

func main() {

	/*
		// Test
		tab := make([]comsoc.Alternative, 10)
		for i, _ := range tab {
			tab[i] = comsoc.Alternative(i + 1)
		}
		fmt.Println(tab)
		fmt.Print(comsoc.Rank(comsoc.Alternative(10), tab))
		mp := make(map[comsoc.Alternative]int)
		for _,k := range tab {
			mp[k] = rand.Intn(10)
		}

		fmt.Println(mp)
		fmt.Println(comsoc.MaxCount(mp))
	*/

	/*
		// TEST CHECK PROFILE
		pref1 := make([]comsoc.Alternative, 11)
		pref2 := make([]comsoc.Alternative, 11)
		pref3 := make([]comsoc.Alternative, 11)
		pref4 := make([]comsoc.Alternative, 11)
		alts := make([]comsoc.Alternative, 11)
		for i, _ := range alts {
			alts[i] = comsoc.Alternative(i + 1)

		}

		for i, _ := range pref1 {
			pref1[i] = comsoc.Alternative(11 - i)

		}
		for i, _ := range pref2 {
			pref2[i] = comsoc.Alternative(i + 1)

		}
		for i, _ := range pref3 {
			pref3[i] = comsoc.Alternative(i + 1)

		}
		for i, _ := range pref4 {
			pref4[i] = comsoc.Alternative(11-i)

		}

		profile := make([][]comsoc.Alternative, 4)
		profile[0] = pref1
		profile[1] = pref2
		profile[2] = pref3
		profile[3] = pref4

		fmt.Println(pref1)
		fmt.Println(pref2)
		fmt.Println(pref3)
		fmt.Println(pref4)
		fmt.Println(comsoc.CheckProfile(profile))
		fmt.Println(comsoc.MajoritySWF(profile))
		fmt.Println(comsoc.MajoritySCF(profile))
		fmt.Println(comsoc.BordaSCF(profile))
	*/

	/*
		prefs := [][]comsoc.Alternative{
			{1, 2, 3},
			{1, 2, 3},
			{3, 2, 1},
			{2, 3, 1},
		}

		res, _ := comsoc.BordaSCF(prefs)

		fmt.Println(res)
	*/
	/*

		prefs := [][]comsoc.Alternative{
			{1, 3, 2},
			{1, 2, 3},
			{2, 3, 1},
			{4, 1, 2},
		}
		thresholds := []int{2, 1, 2, 3}

		res, err := comsoc.ApprovalSCF(prefs, thresholds)

		if err != nil {
			fmt.Println(err)
		}
		if len(res) != 1 || res[0] != 1 {
			fmt.Println("error, 1 should be the only best Alternative")
		}

		fmt.Println(res)
	*/
  /*
	p := [][]comsoc.Alternative{
		{1, 3, 2},
		{1, 2, 3},
		{2, 3, 1},
		{3, 1, 2},
	}

	res, _ := comsoc.STV_SWF(p)

	fmt.Println(res)
  */
  /*

	prefs := [][]comsoc.Alternative{
		{1, 3, 2},
		{1, 2, 3},
		{2, 3, 1},
		{2, 1, 3},
		{3, 1, 2},
		{3, 1, 2},
	}
	//thresholds := []int{2, 1, 2, 3}

	res, err := comsoc.SWFFactory(comsoc.BordaSWF, comsoc.TieBreakFactory([]comsoc.Alternative{3, 2, 1, 4}))(prefs)
	
	fmt.Println(res)
  */
  /*

		prefs := [][]comsoc.Alternative{
			{1, 3, 2},
			{1, 2, 3},
			{2, 3, 1},
			{4, 1, 2},
		}
		thresholds := []int{2, 1, 2, 3}

		res, err := comsoc.ApprovalSCF(prefs, thresholds)

		if err != nil {
			fmt.Println(err)
		}
		if len(res) != 1 || res[0] != 1 {
			fmt.Println("error, 1 should be the only best Alternative")
		}

		fmt.Println(res)
	*/
	/*

		p := [][]comsoc.Alternative{
			{1, 2, 3, 4},
			{1, 2, 3, 4},
			{1, 2, 3, 4},
			{1, 2, 3, 4},
			{1, 2, 3, 4},
			{2, 3, 4, 1},
			{2, 3, 4, 1},
			{2, 3, 4, 1},
			{2, 3, 4, 1},
			{4, 3, 1, 2},
			{4, 3, 1, 2},
			{4, 3, 1, 2},
		}

		fmt.Println(comsoc.CopelandSWF(p))
	*/

}
