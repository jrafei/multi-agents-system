package restserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	rad "gitlab.utc.fr/lagruesy/ia04/demos/restagentdemo"
)

type Alternative int




type RestBallotAgent struct {
	sync.Mutex
	id       string
	addr     string //à vérifier peut etre l'adresse du serveur 
	rule string
	deadline string
	voter_ids []string
	nb_alts int
	tiebreak []Alternative

}

func NewRestBallotAgent(addr string, r string, ru string, d string, vot_ids []string) *RestBallotAgent {
	return &RestBallotAgent{id: addr, addr: addr, rule: r, deadline: d, voter_ids: vot_ids, nb}
}

// Test de la méthode
func (rsa *RestBallotAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (*RestBallotAgent) decodeRequest(r *http.Request) (req rad.Request, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *RestBallotAgent) init_ballot(w http.ResponseWriter, r *http.Request) {
	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()
	//rsa.reqCount++

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// traitement de la requête
	var resp rad.Response

	switch req.Operator {
	case "*":
		resp.Result = req.Args[0] * req.Args[1]
	case "+":
		resp.Result = req.Args[0] + req.Args[1]
	case "-":
		resp.Result = req.Args[0] - req.Args[1]
	default:
		w.WriteHeader(http.StatusNotImplemented)
		msg := fmt.Sprintf("Unkonwn command '%s'", req.Operator)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (rsa *RestBallotAgent) doReqcount(w http.ResponseWriter, r *http.Request) {
	if !rsa.checkMethod("GET", w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)
	rsa.Lock()
	defer rsa.Unlock()
	serial, _ := json.Marshal(rsa.reqCount)
	w.Write(serial)
}

func (rsa *RestBallotAgent) Start() {
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
