package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"ia04/agt"
	rad_t "ia04/agt"
	restserveragent "ia04/agt/restserveragent"
	restvoteragent "ia04/agt/restvoteragent"
	coms "ia04/comsoc"
)

func main() {
	const nVoters = 3 // 3 voters
	const nAlts = 5   // 5 alternatives
	const url1 = ":8080"
	const url2 = "http://localhost:8080"

	server := restserveragent.NewRestServerAgent(url1)

	log.Println("démarrage du serveur...")
	go server.Start()

	//Créer une requete RequestBallot et envoyer vers le serveur
	req := rad_t.RequestBallot{
		Rule:     "majority",
		Deadline: "2023-11-28T23:50:00+02:00",
		Voters:   []string{"ag_id01", "ag_id02", "ag_id03"},
		Nb_alts:  5,
		Tiebreak: []int{4, 2, 3, 5, 1},
	}

	// sérialisation de la requête
	url := url2 + "/new_ballot"
	data, _ := json.Marshal(req) // data de type []octet (json encoding) , traduire la demande en liste de bit (encode)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data)) //resp de type *http.Response , une requete sera envoyé au serveur
	if err != nil {
		log.Printf("[main] erreur %d ...", resp.StatusCode)
		return
	}

	if resp.StatusCode != http.StatusCreated {
		log.Printf("[main] erreur %d ...", resp.StatusCode)
		return
	}

	log.Println("[main] démarrage des voters...")
	votersAgts := make([]restvoteragent.RestVoterAgent, 0, nVoters)
	for i := 0; i < nVoters; i++ {
		id := fmt.Sprintf("ag_id%02d", i+1)
		name := fmt.Sprintf("Voter%02d", i+1)

		var prefs []coms.Alternative
		generated := make(map[int]bool)
		for len(prefs) < nAlts {
			num := rand.Intn(nAlts)

			// Vérifie si l'entier généré est déjà dans la carte
			if !generated[num] {
				generated[num] = true
				prefs = append(prefs, coms.Alternative(num+1))
			}
		}

		ops := make([]int, 1)
		ops[0] = rand.Intn(5) + 1
		agt := restvoteragent.NewRestVoterAgent(id, name, prefs, url2, ops)
		votersAgts = append(votersAgts, *agt)
	}

	//log.Println(votersAgts)

	// A REVOIR QUAND ON ENVOIE UNE REQUETE RESULT TODO
	time.Sleep(10 * time.Second)
	for _, agt := range votersAgts {
		func(agt restvoteragent.RestVoterAgent) {
			go agt.Start("scrutin1")
		}(agt)
	}

	// creation de requete de result
	req_res := agt.RequestVote{
		BallotID: "scrutin1",
	}

	// sérialisation de la requête
	url_res := url2 + "/result"
	data_res, _ := json.Marshal(req_res)

	// envoi de la requête au url
	resp_res, err := http.Post(url_res, "application/json", bytes.NewBuffer(data_res))

	// traitement de la réponse
	if err != nil {
		// A REVOIR [TODO]
		log.Printf("[main] erreur %d ...", resp_res.StatusCode)
		return
	}

	if resp_res.StatusCode != http.StatusOK {
		log.Printf("[main] erreur %d ...", resp_res.StatusCode)
		return
	}

	//log.Println(resp_res)

	fmt.Scanln()
}
