package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

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
		Deadline: "2023-11-04T23:05:08+02:00",
		Voters:   []string{"ag_id01", "ag_id02", "ag_id03"},
		Nb_alts:  5,
		Tiebreak: []int{4, 2, 3, 5, 1},
	}

	// sérialisation de la requête
	url := url2 + "/new_ballot"
	data, _ := json.Marshal(req) // data de type []octet (json encoding) , traduire la demande en liste de bit (encode)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data)) //resp : *http.Response , une requete sera envoyé au serveur
	if err != nil {
		log.Println("erreur 1 ...")
		return
	}

	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		log.Println("erreur 2 ...", resp.StatusCode)
		return
	}

	//log.Println("response : ", resp)

	// Création de plusieurs RequestVote et envoyer au serveur

	// le serveur va renvoyer la requestReponse
	//affichage de reponse du serveur

	log.Println("démarrage des voters...")
	votersAgts := make([]restvoteragent.RestVoterAgent, 0, nVoters)
	for i := 0; i < nVoters; i++ {
		id := fmt.Sprintf("ag_id%02d", i+1)
		name := fmt.Sprintf("Voter02d", i+1)

		var prefs []coms.Alternative
		generated := make(map[int]bool)
		for len(prefs) < nAlts {
			num := rand.Intn(nAlts) // Vous pouvez ajuster la plage selon vos besoins

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

	for _, agt := range votersAgts {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		//for est bcp plus rapide de go , si on met dans for seulement la ligne 40 , on applique le start pour l'agent 99 seulement
		func(agt restvoteragent.RestVoterAgent) {
			go agt.Start("scrutin1")
		}(agt)
	}
	fmt.Scanln()
}
