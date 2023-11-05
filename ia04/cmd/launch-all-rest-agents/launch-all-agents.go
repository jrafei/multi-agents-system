package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	rsa "ia04/agt/restserveragent"
	agent "ia04/agt/restvoteragent"
	comsoc "ia04/comsoc"
)

func main() {
	const nVoters = 5 // 5 voters
	const nAlts = 5   // 5 alternatives
	const port = ":8080"
	const url_server = "http://localhost:8080"

	server := rsa.NewRestServerAgent(port)

	log.Println("démarrage du serveur...")
	go server.Start()

	//Créer un agent qui demande la création d'un ballot
	ballot_agent := agent.NewAgent("ballot_agt1", "ballot_agt", nil, nil)

	deadline := time.Now().Add(time.Second * 10).Format(time.RFC3339)
	Rule := "majority"
	Voters := []string{"ag_id01", "ag_id02", "ag_id03", "ballot_agt1"}
	Nb_alts := 5
	Tiebreak := []int{4, 2, 3, 5, 1}
	// demande de creation d'un ballot
	ballot_agent.CreateBallot(Rule, deadline, Voters, Nb_alts, Tiebreak, url_server)

	log.Println("[main] démarrage des voters...")
	votersAgts := make([]agent.Agent, 0, nVoters)
	for i := 0; i < nVoters; i++ {
		id := fmt.Sprintf("ag_id%02d", i+1)
		name := fmt.Sprintf("Voter%02d", i+1)

		var prefs []comsoc.Alternative
		generated := make(map[int]bool)
		for len(prefs) < nAlts {
			num := rand.Intn(nAlts)

			// Vérifie si l'entier généré est déjà dans la carte
			if !generated[num] {
				generated[num] = true
				prefs = append(prefs, comsoc.Alternative(num+1))
			}
		}

		ops := make([]int, 1)
		ops[0] = rand.Intn(5) + 1
		agt := agent.NewAgent(id, name, prefs, ops)
		votersAgts = append(votersAgts, *agt)

		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		//for est bcp plus rapide de go , si on met dans for seulement la ligne 40 , on applique le start pour l'agent 99 seulement
		func() {
			go agt.Vote("scrutin1", url_server)
		}()
	}

	go ballot_agent.Vote("scrutin1", url_server)

	for {
		// Récupération du résultat du scrutin
		if time.Now().Format(time.RFC3339) > deadline {
			winner, ranking, err := votersAgts[rand.Intn(nVoters)].GetResult("scrutin1", url_server)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("winner : ", winner)
				fmt.Println("ranking : ", ranking)
			}

			return
		}
	}
}
