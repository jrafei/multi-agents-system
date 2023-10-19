package restserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"


	rad "ia04/agt/restballotagent"
	rad_t "ia04/agt"
	comsoc "ia04/comsoc"
)

type RestServerAgent struct {
	sync.Mutex
	id       string
	reqCount int
	addr     string
	ballots  map[string]chan string // associe ballot-id et chan associé pour communiquer avec le serveur
	channel  chan string
}

func NewRestServerAgent(addr string) *RestServerAgent {
	return &RestServerAgent{id: addr, addr: addr, ballots: make(map[string]chan string), channel: make(chan string)}
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
func (*RestServerAgent) decodeRequestBallot(r *http.Request) (req rad_t.RequestBallot, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

// Execution de la création d'un ballot
func (rsa *RestServerAgent) init_ballot(w http.ResponseWriter, r *http.Request) {
	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()
	rsa.reqCount++

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequestBallot(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// traitement de la requête
	var resp rad_t.Response

	fmt.Println(req)
	
	// Vérification des paramètres
	if req.Nb_alts < 0 {

		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Negative number of alternatives (%s)", req.Nb_alts)
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
		msg := fmt.Sprintf("the deadline has already been passed (%s)", req.Deadline)
		w.Write([]byte(msg))
		return
	}

	// Vérification du tiebreak
	alts := make([]comsoc.Alternative, req.Nb_alts)
	tieb := make([]comsoc.Alternative, req.Nb_alts)
	for i, _ := range tieb {
		alts[i] = comsoc.Alternative(i + 1)
		tieb[i] = comsoc.Alternative(req.Tiebreak[i])
	}
	fmt.Println(alts)
	fmt.Println(tieb)
	if comsoc.CheckProfile(tieb, alts) != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("The tiebreak doesn't correctly represent the alternatives")
		w.Write([]byte(msg))
		return
	}
	// faire vérif méthodes

	// création d'une ballot si tout est conforme
	ballot_id := string("scrutin" + strconv.Itoa(len(rsa.ballots)+1))
	ballot_ch := make(chan string)
	rsa.ballots[ballot_id] = ballot_ch
	ballot := rad.NewRestBallotAgent(ballot_id, req.Rule, req.Deadline, req.Voters, req.Nb_alts, tieb, ballot_ch)

	// Lancement de la ballot par une go routine (ajout)
	go ballot.Start(ballot_ch)

	/*
		switch req.Rule {
		case "majority":
			resp.Result = req.Args[0] * req.Args[1]
		default:
			w.WriteHeader(http.StatusNotImplemented)
			msg := fmt.Sprintf("Unkonwn command '%s'", req.Operator)
			w.Write([]byte(msg))
			return
		}
	*/
	// à modfier
	resp.Ballot_id = ballot_id
	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (rsa *RestServerAgent) doReqcount(w http.ResponseWriter, r *http.Request) {
	if !rsa.checkMethod("GET", w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)
	rsa.Lock()
	defer rsa.Unlock()
	serial, _ := json.Marshal(rsa.reqCount)
	w.Write(serial)
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.init_ballot)
	mux.HandleFunc("/vote", rsa.doReqcount)
	mux.HandleFunc("/result", rsa.doReqcount)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}
