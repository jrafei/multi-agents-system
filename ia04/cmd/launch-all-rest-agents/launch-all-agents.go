package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"strconv"
	"strings"
	"time"
	rad_t "ia04/agt"
	restserveragent "ia04/agt/restserveragent"
	restvoteragent "ia04/agt/restvoteragent"
	coms "ia04/comsoc"
)

func main() {
	const nVoters = 5 // 5 voters
	const nAlts = 5   // 5 alternatives
	const url1 = ":8080"
	const url2 = "http://localhost:8080"

	server := restserveragent.NewRestServerAgent(url1)

	log.Println("démarrage du serveur...")
	go server.Start()

	//Créer une requete RequestBallot et envoyer vers le serveur
	deadline := time.Now().Add(time.Second*10).Format(time.RFC3339)
	req := rad_t.RequestBallot{
		Rule:     "majority",

		Deadline: deadline,			// On implémente une deadline à + 10 secondes

		Voters:   []string{"ag_id01", "ag_id02", "ag_id03"},
		Nb_alts:  5,
		Tiebreak: []int{4, 2, 3, 5, 1},
	}

	// sérialisation de la requête
	url := url2 + "/new_ballot"
	data, _ := json.Marshal(req) // data de type []octet (json encoding) , traduire la demande en liste de bit (encode)

	// envoi de la requête

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data)) //resp : *http.Response , une requete sera envoyé au serveur
	/*
	if err != nil {
		log.Println("erreur 1 ...")
		//return
	}
	*/
	/*
	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		log.Println("erreur 2 ...", resp.StatusCode)
		//return
	}
	*/
	if err != nil {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		log.Println("erreur ", resp.StatusCode)
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

		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		//for est bcp plus rapide de go , si on met dans for seulement la ligne 40 , on applique le start pour l'agent 99 seulement
		func() {
			go agt.Start("scrutin1")
		}()
	}

	//log.Println(votersAgts)


	/*
	for _, agt := range votersAgts {
		func(agt restvoteragent.RestVoterAgent) {
			go agt.Start("scrutin1")
		}(agt)
	}
	*/

	for{
		// Récupération du résultat du scrutin
		if time.Now().Format(time.RFC3339) > deadline {
			resp_s,err := votersAgts[rand.Intn(nVoters)].DoRequestResult("scrutin1")
			if err != nil {
				log.Println("[CLIENT] An error occured : " + err.Error())
			}
			if resp_s.Status != http.StatusOK{
				log.Println(http.StatusText(resp_s.Status))
				
			}else {
				ranking := make([]string,len(resp_s.Ranking))
				for i,v := range resp_s.Ranking{
					ranking[i] = strconv.Itoa(v)
				}
				log.Println("[CLIENT] Client received (" + strconv.Itoa(resp_s.Status) + ") : " + "\nWinner : " +  strconv.Itoa(resp_s.Winner) + "\nRanking : " + strings.Join(ranking, ","))
			}
			return

		}
	}
}