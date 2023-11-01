package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	rsa "ia04/agt/restserveragent"
	agent "ia04/agt/restvoteragent"
	utils "ia04/agt/utils"
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

	log.Println("test")
	list_ballot := [6]string{"majority", "borda", "approval", "condorcet", "copeland", "stv"}
	m := make(map[int]string) // map , clé = nom de scrutin, valeur = deadline
	//var list_deadline []string
	for ind, ballot := range list_ballot {
		//Créer une requete RequestBallot et envoyer vers le serveur
		mult := (ind + 1) * 5
		deadline := time.Now().Add(time.Second * time.Duration(mult)).Format(time.RFC3339)
		log.Println("=========================================" + deadline)
		m[ind+1] = deadline // false pour résultat non publié encore
		//list_deadline = append(list_deadline, deadline)
		req := utils.RequestBallot{
			Rule:     ballot,
			Deadline: deadline, // On implémente une deadline à + 10 secondes
			Voters:   []string{"ag_id01", "ag_id02", "ag_id03"},
			Nb_alts:  5,
			Tiebreak: []int{4, 2, 3, 5, 1},
		}

		// sérialisation de la requête
		url_request := url_server + "/new_ballot"
		data, _ := json.Marshal(req) // data de type []octet (json encoding) , traduire la demande en liste de bit

		// envoi de la requête
		resp, err := http.Post(url_request, "application/json", bytes.NewBuffer(data)) //resp : *http.Response , une requete sera envoyé au serveur

		if err != nil {
			log.Println("erreur status code : ", resp.StatusCode)
			log.Println("erreur status 		: ", resp.Status)
			return
		}
	}

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

		func() {
			go agt.Vote("scrutin1", url_server)
			go agt.Vote("scrutin2", url_server)
			go agt.Vote("scrutin3", url_server)
			go agt.Vote("scrutin4", url_server)
			go agt.Vote("scrutin5", url_server)
			go agt.Vote("scrutin6", url_server)
		}()
	}

	for {
		for i, deadl := range m {
		
			if time.Now().Format(time.RFC3339) > deadl {
				log.Println("============================== Result " + strconv.Itoa(i) + " ========================")
				winner, ranking, err := votersAgts[rand.Intn(nVoters)].GetResult("scrutin"+strconv.Itoa(i), url_server)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("winner : ", winner)
					fmt.Println("ranking : ", ranking)
				}
				delete(m, i)
			}
		}
		if len(m) == 0 {
			return
		}
	}
}
