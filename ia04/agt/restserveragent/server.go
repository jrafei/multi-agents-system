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
	request "ia04/agt/request"
	comsoc "ia04/comsoc"
)

type RestServerAgent struct {
	sync.Mutex
	id      string
	addr    string
	ballots map[string]chan request.RequestVoteBallot // associe ballot-id et chan associé pour communiquer avec le serveur
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
	return &RestServerAgent{id: addr, addr: addr, ballots: make(map[string]chan request.RequestVoteBallot)}
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
	    - booléen : vrai si identique, faux sinon

======================================
*/
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp_finale := request.Response{Info: "Mauvaise méthode HTTP réceptionnée."}
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
func (*RestServerAgent) decodeRequestBallot(r *http.Request) (req request.RequestBallot, err error) {
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
func (*RestServerAgent) decodeRequestVote(r *http.Request) (req request.RequestVote, err error) {
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

	// On lock le système pour ne pas avoir de conflit
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête http -> initialisation de structure RequestBallot 'req'
	req, err := rsa.decodeRequestBallot(r)
	log.Println("[SERVER] Received ballot creation request : ", req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp_finale := request.Response{Info: "Impossible de comprendre la requête : " + err.Error()}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return
	}

	// traitement de la requête
	var resp request.Response

	// Vérification de la méthode de vote
	switch req.Rule {
	case "majority", "borda", "approval", "stv", "copeland", "condorcet", "kramer-simpson", "young-condorcet", "dodgson", "kemeny-young":
		break
	default:
		w.WriteHeader(http.StatusNotImplemented)
		resp_finale := request.Response{Info: "Méthode de vote inconnue."}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return
	}

	// Vérification des paramètres
	deadline, time_err := time.Parse(time.RFC3339, req.Deadline)
	time_now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	if req.Nb_alts <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		resp_finale := request.Response{Info: "Nombre nul ou négatif d'alternatives."}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return
	} else if len(req.Voters) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		resp_finale := request.Response{Info: "Aucun voteur sur la liste."}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return
	} else if time_err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp_finale := request.Response{Info: "Le format de la deadline n'est pas reconnu (attendu : RFC3339)."}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return
	} else if time_now.After(deadline) {
		w.WriteHeader(http.StatusBadRequest)
		resp_finale := request.Response{Info: "La deadline est déjà dépassée."}
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
		resp_finale := request.Response{Info: "Le tiebreak ne représente pas correctement les alternatives."}
		serial, _ := json.Marshal(resp_finale)
		w.Write([]byte(serial))
		return
	}

	// création d'une ballot si tout est conforme
	ballot_id := string("scrutin" + strconv.Itoa(len(rsa.ballots)+1))
	ballot_ch := make(chan request.RequestVoteBallot)
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
			- 'w' : http ResponseWriter pour réponse
			- 'r' : requete http contenant les caractéristiques du ballot

======================================
*/
func (rsa *RestServerAgent) ballotHandler(action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// On lock le système pour ne pas avoir de conflit
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
				resp_finale := request.Response{Info: "Impossible de comprendre la requête : " + err.Error()}
				serial, _ := json.Marshal(resp_finale)
				w.Write([]byte(serial))
				return
			}

			// traitement de la requête
			var resp request.RequestVoteBallot

			// Vérification du BallotID
			ballot_chan, exists := rsa.ballots[req.BallotID]
			if !exists {
				if action == "vote" {
					w.WriteHeader(http.StatusBadRequest)
				} else if action == "result" {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
				resp_finale := request.Response{Info: "Le ballotID n'est pas reconnu."}
				serial, _ := json.Marshal(resp_finale)
				w.Write([]byte(serial))
				return
			}

			vote_req := request.RequestVoteBallot{RequestVote: &req, Action: action, StatusCode: 0, Msg: ""}

			// Transmission de la requête au ballot correspondant
			ballot_chan <- vote_req
			// Attente de la response du ballot
			resp = <-ballot_chan
			// Transmission de la réponse du ballot au client
			switch action {
			case "vote":
				w.WriteHeader(resp.StatusCode)
				resp_finale := request.Response{Info: resp.Msg}
				serial, _ := json.Marshal(resp_finale)
				w.Write([]byte(serial))
			case "result":
				if resp.StatusCode == http.StatusOK {
					w.WriteHeader(http.StatusOK)
					resp_finale := request.Response{Winner: resp.Winner, Ranking: resp.Ranking, Info: resp.Msg}
					serial, _ := json.Marshal(resp_finale)
					w.Write([]byte(serial))
				} else {
					w.WriteHeader(resp.StatusCode)
					resp_finale := request.Response{Info: resp.Msg}
					serial, _ := json.Marshal(resp_finale)
					w.Write([]byte(serial))
				}
			}
		case http.MethodGet:
			// Répondre à une requête GET
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("GET method not implemented"))
		case http.MethodPut:
			// Répondre à une requête PUT
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("PUT method not implemented"))
		case http.MethodDelete:
			// Répondre à une requête DELETE
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("DELETE method not implemented"))
		default:
			// Répondre pour toute autre méthode HTTP non prise en charge
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
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
