package restserveragent

/*
 * TOCHECK : ERROR 501 /vote (à quoi cela correspond-il ?)
 */
import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	ballotagent "ia04/agt/ballotagent"
	utils "ia04/agt/utils"
	comsoc "ia04/comsoc"
)

type RestServerAgent struct {
	sync.Mutex
	id      string
	addr    string
	ballots map[string]chan utils.RequestVoteBallot // associe ballot-id et chan associé pour communiquer avec le serveur
}

/*
======================================

	  @brief :
	  'Constructeur de la classe.'
	  @params :
		- 'addr' : url/port du serveur
	  @returned :
	    -  Un pointeur sur le serveur créé.

======================================
*/
func NewRestServerAgent(addr string) *RestServerAgent {
	return &RestServerAgent{id: addr, addr: addr, ballots: make(map[string]chan utils.RequestVoteBallot)}
}

/*
======================================

	  @brief :
	  'Méthode de vérification de la bon type de requete http attendu (GET,POST...)'
	  @params :
		- 'method' : type de requete attendu
		- 'w' : http ResponseWriter pour réponse
		- 'r' : requete http à vérifier
	  @returned :
	    - booléen

======================================
*/
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp_finale := utils.Response{Info: "Mauvaise méthode HTTP réceptionnée."}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return false
	}
	return true
}

/*
======================================

	  @brief :
	  'Méthode pour décodage une requête de création de ballot '
	  @params :
		- 'r' : requete http contenant les caractéristiques du ballot (binaire)
	  @returned :
	    - 'req' : structure RequestVote contenant les informations traduites envoyées par le client
		- 'err' : variable d erreur

======================================
*/
func (*RestServerAgent) decodeRequestBallot(r *http.Request) (req utils.RequestBallot, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

/*
======================================

	  @brief :
	  'Méthode pour décodage une requête de proposition de vote '
	  @params :
		- 'r' : requete http contenant les caractéristiques du vote (binaire)
	  @returned :
	    - 'req' : structure RequestVote contenant les informations traduites envoyées par le client
		- 'err' : variable d erreur

======================================
*/
func (*RestServerAgent) decodeRequestVote(r *http.Request) (req utils.RequestVote, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

/*
======================================

	  @brief :
	  'Handler de création d un ballot.'
	  @params :
		- 'w' : http ResponseWriter pour réponse
		- 'r' : requete http contenant les caractéristiques du ballot

======================================
*/
func (rsa *RestServerAgent) init_ballot(w http.ResponseWriter, r *http.Request) {

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
		resp_finale := utils.Response{Info: "Impossible de comprendre la requête : " + err.Error()}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return
	}

	// traitement de la requête
	var resp utils.Response

	// fmt.Println("Request /init_ballot :")
	// fmt.Println(req)

	// Vérification de la méthode de vote
	switch req.Rule {
	case "majority", "borda", "approval", "stv", "copeland", "condorcet":
		break
	default:
		w.WriteHeader(http.StatusNotImplemented)
		resp_finale := utils.Response{Info: "Méthode de vote inconnue."}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return
	}

	// Vérification des paramètres
	if req.Nb_alts < 0 {
		w.WriteHeader(http.StatusBadRequest)
		resp_finale := utils.Response{Info: "Nombre négatif d'alternatives."}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return
	} else if len(req.Voters) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		resp_finale := utils.Response{Info: "Nombre négatif de voteurs."}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return
	} else if req.Deadline <= time.Now().Format(time.RFC3339) {
		w.WriteHeader(http.StatusBadRequest)
		resp_finale := utils.Response{Info: "La deadline est déjà dépassée."}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
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
		resp_finale := utils.Response{Info: "Le tiebreak ne représente pas correctement les alternatives."}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return
	}

	// création d'une ballot si tout est conforme
	ballot_id := string("scrutin" + strconv.Itoa(len(rsa.ballots)+1))
	ballot_ch := make(chan utils.RequestVoteBallot)
	rsa.ballots[ballot_id] = ballot_ch
	ballot := ballotagent.NewRestBallotAgent(ballot_id, req.Rule, req.Deadline, req.Voters, req.Nb_alts, tieb, ballot_ch)

	// Lancement de la ballot par une go routine (ajout)
	go ballot.Start()

	// Initialisation de la réponse
	resp.Ballot_id = ballot_id
	resp.Info = "Ballot créé."
	w.WriteHeader(http.StatusCreated)
	serial, _ := json.Marshal(resp)
	w.Write(serial)

	// /********DEBUG********/
	// fmt.Println("-----------------")
	// fmt.Println("Updated server after /init_ballot :")
	// fmt.Println(rsa.id)
	// fmt.Println(rsa.addr)
	// fmt.Println(rsa.ballots)
	// fmt.Println("-----------------")
	// /*********************/

}

/*
======================================

	  @brief :
	  'Méthode de déclenchement de l action souhaitée par le client.'
	  Utilisation d un wrapper pour ajouter le paramètre action.
	  Permet d éviter la duplication de code inutile aux différents types d action.
	  @params :
		- 'action' : Le procédé souhaité
	  @returned :
	    - fonction associée à l action désirée

======================================
*/
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
			resp_finale := utils.Response{Info: "Impossible de comprendre la requête : " + err.Error()}
			serial, _ := json.Marshal(resp_finale)
			w.Write([]byte(serial))
			return
		}

		// traitement de la requête
		var resp utils.RequestVoteBallot

		// /********DEBUG********/
		// fmt.Printf("[%s] Request /%s from client to server :\n", req.AgentID, action)
		// fmt.Println("RequestVote : ", req)

		// Vérification du BallotID
		ballot_chan, exists := rsa.ballots[req.BallotID]
		if !exists {
			if action == "vote" { //TODO : à vérifier
				w.WriteHeader(http.StatusBadRequest)
			} else if action == "result" {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			resp_finale := utils.Response{Info: "Le ballotID n'est pas reconnu."}
			serial, _ := json.Marshal(resp_finale)
			w.Write([]byte(serial))
			return
		}

		vote_req := utils.RequestVoteBallot{RequestVote: &req, Action: action, StatusCode: 0, Msg: ""}

		// /********DEBUG********/
		// fmt.Printf("[%s] Request /%s from server to ballot :\n", req.AgentID, action)
		// fmt.Println("RequestVoteBallot : ", vote_req)
		// /*********************/

		// Transmission de la requête au ballot correspondant
		ballot_chan <- vote_req
		// Attente de la response du ballot
		resp = <-ballot_chan

		// /********DEBUG********/
		// fmt.Println("Reponse from server to client : ")
		// fmt.Printf("[%s] Action :%s \n", req.AgentID, resp.Action)
		// fmt.Printf("[%s] Status Code : %d \n", req.AgentID, resp.StatusCode)
		// fmt.Printf("[%s] Msg : %s \n", req.AgentID, resp.Msg)
		// fmt.Printf("[%s] Winner : %d\n", req.AgentID, resp.Winner)
		// fmt.Printf("[%s] Ranking : %d \n", req.AgentID, resp.Ranking)
		// /*********************/

		// Transmission de la réponse du ballot au client
		switch action {
		case "vote":
			w.WriteHeader(resp.StatusCode)
			resp_finale := utils.Response{Info: resp.Msg}
			serial, _ := json.Marshal(resp_finale)
			w.Write([]byte(serial))
		case "result":
			if resp.StatusCode == http.StatusOK {
				w.WriteHeader(http.StatusOK)
				resp_finale := utils.Response{Winner: resp.Winner, Ranking: resp.Ranking, Info: resp.Msg}
				serial, _ := json.Marshal(resp_finale)
				w.Write([]byte(serial))
			} else {
				w.WriteHeader(resp.StatusCode)
				resp_finale := utils.Response{Info: resp.Msg}
				serial, _ := json.Marshal(resp_finale)
				w.Write([]byte(serial))
			}

		}
	}
}

/*
======================================

	@brief :
	'Procédure de mise en fonction du serveur. Elle crée et écoute les requêtes http entrantes.'

======================================
*/
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
