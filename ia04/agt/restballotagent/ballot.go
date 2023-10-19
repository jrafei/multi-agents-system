package restserveragent

import (
	"sync"
	com "ia04/comsoc"
)



type RestBallotAgent struct {
	sync.Mutex
	id          string
	rule        string
	deadline    string
	voter_ids   []string
	nb_alts     int
	tiebreak    []com.Alternative
	server_chan chan string
}

func NewRestBallotAgent(i string, ru string, d string, vot_ids []string, alts int, tieb []com.Alternative, ch chan string) *RestBallotAgent {
	return &RestBallotAgent{id: i, rule: ru, deadline: d, voter_ids: vot_ids, nb_alts: alts, tiebreak: tieb, server_chan: ch}
}

/*
// Test de la méthode
func (rsa *RestBallotAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}
*/

/*
func (*RestBallotAgent) decodeRequest(r *http.Request) (req rad.Request, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}
*/
/*
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
*/

func (rsa *RestBallotAgent) Start(chan string) {
	// si le channel reçoit une demande, on lace la méthode associée
	for {

	}
}
