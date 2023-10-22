package main

import (
	"fmt"
	"log"
	"math/rand"

	restserveragent "ia04/agt/restserveragent"

	restclientagent "ia04/agt/restvoteragent"
	coms "ia04/comsoc"
)

func main() {
	const n = 100
	const url1 = ":8080"
	const url2 = "http://localhost:8080"

	votersAgts := make([]restclientagent.RestVoterAgent, 0, n)
	server := restserveragent.NewRestServerAgent("url1")

	log.Println("démarrage du serveur...")
	go server.Start()

	// initialisation d'un ballot ??? envoie d'une requete post localhost:8080/new_ballot
	//Créer un requete new_ballot

	//log.Println("démarrage des voters...")
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("ag_id02d", i)
		name := fmt.Sprintf("Voter02d", i)

		// le test est pour 5 alternatives
		prefs := make([]coms.Alternative, 5)
		for ind := 0; ind < 10; ind++ {
			prefs[ind] = coms.Alternative(rand.Intn(5) + 1)
		}
		agt := restclientagent.NewRestVoterAgent(id, name, prefs, url2)
		votersAgts = append(votersAgts, *agt)
	}

	for _, voter := range votersAgts {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt restclientagent.RestVoterAgent) {
			go voter.Start()
		}(voter)
	}

	fmt.Scanln()
}
