package restserveragent

import (
	rad_t "ia04/agt"
	"ia04/comsoc"
	com "ia04/comsoc"
	"sync"
	"time"
)

type RestBallotAgent struct {
	sync.Mutex
	id          string
	rule        string
	deadline    string
	voter_ids   map[string]bool // Le booléen permet de savoir si l'agent a déjà voté
	nb_alts     int
	tiebreak    []com.Alternative
	server_chan chan rad_t.RequestVoteBallot
	profile     com.Profile
	options     [][]int
}

func NewRestBallotAgent(i string, ru string, d string, vot_ids []string, alts int, tieb []com.Alternative, ch chan rad_t.RequestVoteBallot) *RestBallotAgent {
	voters := make(map[string]bool, 0)
	for _, id := range vot_ids {
		// On signale qu'aucun voteur n'a encore voté
		voters[id] = false
	}
	return &RestBallotAgent{id: i, rule: ru, deadline: d, voter_ids: voters, nb_alts: alts, tiebreak: tieb, server_chan: ch, profile: make(com.Profile, 0), options: make([][]int, 0)}
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

func (rsa *RestBallotAgent) Start(chan rad_t.RequestVoteBallot) {
	// si le channel reçoit une demande, on lace la méthode associée
	for {

		req := <-rsa.server_chan
		// Selection de l'action à effectuer
		switch req.Action {
		case "vote":
			resp := rsa.Vote(req)
			// Transmission de la response au serveur
			rsa.server_chan <- resp
		case "result":
			rsa.Result()
		}

	}
}

func (rsa *RestBallotAgent) Vote(vote rad_t.RequestVoteBallot) (resp rad_t.RequestVoteBallot) {

	// Vérification de la deadline
	if rsa.deadline <= time.Now().Format(time.RFC3339) {
		resp.StatusCode = 503
		resp.Msg = "la deadline est dépassée"
		return
	}

	// Vérification de l'AgentID
	has_voted, exists := rsa.voter_ids[vote.AgentID]
	if !exists {
		resp.StatusCode = 400
		resp.Msg = "bad request, le voteur n'est pas sur la liste"
		return
	}
	if has_voted {
		resp.StatusCode = 403
		resp.Msg = "vote déjà effectué"
		return
	}

	// Vérification des préférences
	prefs := make([]comsoc.Alternative, len(vote.Preferences))
	alts := make([]comsoc.Alternative, rsa.nb_alts)
	for i,_ := range alts {
		alts[i] = comsoc.Alternative(i + 1)
	}
	for i,_ := range prefs {
		prefs[i] = comsoc.Alternative(vote.Preferences[i])
	}

	if comsoc.CheckProfile(prefs, alts) != nil {
		resp.StatusCode = 400
		resp.Msg = "bad request, les préférences ne sont pas conformes"
		return
	}

	// Ajout des préférences dans le profil
	rsa.profile = append(rsa.profile, prefs)

	// Ajout des options
	if vote.Options != nil {
		rsa.options = append(rsa.options, vote.Options)
	}

	rsa.voter_ids[vote.AgentID] = true // on indique que l'agent a voté
	resp.StatusCode = 200
	resp.Msg = "vote pris en compte"
	return
}

func (rsa *RestBallotAgent) Result() (resp rad_t.RequestVoteBallot) {
	return
}
