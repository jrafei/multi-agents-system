package restserveragent

/*
 * TOCHECK : ERROR 501 /vote (à quoi cela correspond-il ?)
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	rad_t "ia04/agt"
	rad "ia04/agt/restballotagent"
	comsoc "ia04/comsoc"
)

type RestServerAgent struct {
	sync.Mutex
	id      string
	addr    string
	ballots map[string]chan rad_t.RequestVoteBallot // associe ballot-id et chan associé pour communiquer avec le serveur
	// channel chan rad_t.RequestVoteBallot            //Sert à quoi ???
}

func NewRestServerAgent(addr string) *RestServerAgent {
	// return &RestServerAgent{id: addr, addr: addr, ballots: make(map[string]chan rad_t.RequestVoteBallot), channel: make(chan rad_t.RequestVoteBallot)}
	return &RestServerAgent{id: addr, addr: addr, ballots: make(map[string]chan rad_t.RequestVoteBallot)}

}

// Test de la méthode
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

// Décode une requête de creation d'un ballot
// Renvoir la structure RequestBallot contenant les infos envoyé par le client et une erreur si existe
func (*RestServerAgent) decodeRequestBallot(r *http.Request) (req rad_t.RequestBallot, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

// Décode une requête de proposition de vote
// Renvoie la structure RequestVote contenant les informations envoyé par le client et une erreur si existe
func (*RestServerAgent) decodeRequestVote(r *http.Request) (req rad_t.RequestVote, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

// Handler pour la création d'un ballot
// prend en argument :
// 1 - http ResponseWriter
// 2 - la requete http contenant les infos de ballot à créer codé par Marshal
// création du ballot si les infos sont conformes, et lance le ballot crée
func (rsa *RestServerAgent) init_ballot(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().Format(time.RFC3339))
	// On lock le système pour ne pas avoir de conflit (TODO : à modifier peut-être)
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête http -> initialisation de structure RequestBallot 'req'
	req, err := rsa.decodeRequestBallot(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// traitement de la requête
	var resp rad_t.Response

	fmt.Println("-----------------")
	fmt.Println("[DBG] Request /init_ballot :")
	fmt.Println(req)
	fmt.Println("-----------------")

	// Vérification des paramètres
	if req.Nb_alts < 0 {

		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Negative number of alternatives (%d)", req.Nb_alts)
		w.Write([]byte(msg))
		return
	} else if len(req.Voters) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Negative number of voters (%d)", len(req.Voters))
		w.Write([]byte(msg))
		return
	} else if len(req.Tiebreak) != req.Nb_alts {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("The tiebreak doesn't represent correctly the alternatives (TB : %d - #alts : %d )", len(req.Voters), req.Nb_alts)
		w.Write([]byte(msg))
		return
	} else if req.Deadline <= time.Now().Format(time.RFC3339) {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("The deadline has already been passed (%s)", req.Deadline)
		w.Write([]byte(msg))
		return
	}

	// Vérification de la méthode de vote
	switch req.Rule {
	case "majority", "borda", "approval", "stv", "copeland", "condorcet":
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Unknown rule (%s)", req.Rule)
		w.Write([]byte(msg))
		return
	}

	// Vérification du tiebreak : verifie si toutes les alternatives apparaissent dans la liste tiebreak

	tieb := make([]comsoc.Alternative, len(req.Tiebreak))
	alts := make([]comsoc.Alternative, req.Nb_alts)
	for i, _ := range alts {
		alts[i] = comsoc.Alternative(i + 1)
	}
	for i, _ := range tieb {
		tieb[i] = comsoc.Alternative(req.Tiebreak[i])
	}

	if comsoc.CheckProfile(tieb, alts) != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "The tiebreak doesn't correctly represent the alternatives"
		w.Write([]byte(msg))
		return
	}

	// création d'une ballot si tout est conforme
	ballot_id := string("scrutin" + strconv.Itoa(len(rsa.ballots)+1))
	ballot_ch := make(chan rad_t.RequestVoteBallot)
	rsa.ballots[ballot_id] = ballot_ch
	ballot := rad.NewRestBallotAgent(ballot_id, req.Rule, req.Deadline, req.Voters, req.Nb_alts, tieb, ballot_ch)

	// Lancement de la ballot par une go routine (ajout)
	go ballot.Start(ballot_ch)

	resp.Ballot_id = ballot_id
	w.WriteHeader(http.StatusCreated)
	serial, _ := json.Marshal(resp)
	w.Write(serial)

	/********DEBUG********/
	fmt.Println("-----------------")
	fmt.Println("[DBG] Updated server after /init_ballot :")
	fmt.Println(rsa.id)
	fmt.Println(rsa.addr)
	fmt.Println(rsa.ballots)
	fmt.Println("-----------------")
	/*********************/
}

// Utilisation d'un wrapper pour ajouter le paramètre action (type d'action à effectuer pour le ballot, ex : vote, result)
// Permet d'éviter la duplication de code inutile aux deux types d'action
func (rsa *RestServerAgent) ballotHandler(action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// On lock le système pour ne pas avoir de conflit (TODO : à modifier peut-être)
		rsa.Lock()
		defer rsa.Unlock()

		// vérification de la méthode de la requête
		if !rsa.checkMethod("POST", w, r) {
			return
		}

		// décodage de la requête
		req, err := rsa.decodeRequestVote(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}

		// traitement de la requête
		var resp rad_t.RequestVoteBallot

		/********DEBUG********/
		fmt.Println("-----------------")
		fmt.Printf("[DBG] [%s] Request /%s from client to server :\n", req.AgentID, action)
		fmt.Println("RequestVote : ", req)
		fmt.Println("-----------------")
		/*********************/

		// Vérification du BallotID
		ballot_chan, exists := rsa.ballots[req.BallotID]
		if !exists {
			if action == "vote" {
				w.WriteHeader(http.StatusBadRequest)
			} else if action == "result" {
				w.WriteHeader(404)
			}
			msg := fmt.Sprintf("not found, the ballot '%s' doesn't exist", req.BallotID)
			w.Write([]byte(msg))
			return
		}

		vote_req := rad_t.RequestVoteBallot{RequestVote: &req, Action: action, StatusCode: 0, Msg: ""}

		/********DEBUG********/
		fmt.Println("-----------------")
		fmt.Printf("[DBG] [%s] Request /%s from server to ballot :\n", req.AgentID, action)
		fmt.Println("RequestVoteBallot : ", vote_req)
		fmt.Println("-----------------")
		/*********************/

		// Transmission de la requête au ballot correspondant
		ballot_chan <- vote_req
		// Attente de la response du ballot
		resp = <-ballot_chan

		/********DEBUG********/
		fmt.Println("Reponse from server to client : ")
		fmt.Printf("[%s] Action :%s \n", req.AgentID, resp.Action)
		fmt.Printf("[%s] Status Code : %d \n", req.AgentID, resp.StatusCode)
		fmt.Printf("[%s] Msg : %s \n", req.AgentID, resp.Msg)
		fmt.Printf("[%s] Winner : %d\n", req.AgentID, resp.Winner)
		fmt.Printf("[%s] Ranking : %d \n", req.AgentID, resp.Ranking)
		/*********************/

		// Transmission de la réponse du ballot au client
		switch action {
		case "vote":
			w.WriteHeader(resp.StatusCode)
			msg := resp.Msg
			w.Write([]byte(msg))
		case "result":
			if resp.StatusCode == 200 {
				w.WriteHeader(http.StatusOK)
				resp_finale := rad_t.Response{Winner: resp.Winner, Ranking: resp.Ranking}
				serial, _ := json.Marshal(resp_finale)
				w.Write([]byte(serial))
			} else {
				w.WriteHeader(resp.StatusCode)
				msg := resp.Msg
				w.Write([]byte(msg))
			}

		}
	}
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.init_ballot)
	mux.HandleFunc("/vote", rsa.ballotHandler("vote"))
	mux.HandleFunc("/result", rsa.ballotHandler("result"))

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr, //adresse de localhost
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}
