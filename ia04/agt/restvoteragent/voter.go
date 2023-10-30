package restvoteragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"ia04/agt"
	rad_t "ia04/agt"

	//rad "ia04/agt/restballotagent"
	//agt "ia04/agt/agent"
	coms "ia04/comsoc"
)

type RestVoterAgent struct {
	agt  *rad_t.Agent
	url_server  string //localhost:8080
	opts []int
}

func NewRestVoterAgent(id string, n string, p []coms.Alternative, u string, op []int) *RestVoterAgent {
	ag := rad_t.NewAgent(id, n, p)
	return &RestVoterAgent{ag, u, op}
}


// Décode une réponse
// Renvoie la structure Response, la réponse du server au client
func (*RestVoterAgent) decodeResponse(r *http.Response) (rep rad_t.Response, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &rep)
	rep.Status = r.StatusCode
	return
}

/*
renvoie la réponse du serveur ou une erreur
*/
func (rva *RestVoterAgent) doRequestVoter(ballotID string) (res rad_t.Response, err error) {
	// creation de requete de vote
	req := agt.RequestVote{
		AgentID:     string(rva.agt.ID),
		BallotID:    ballotID,
		Preferences: rva.agt.Prefs,
		Options:     rva.opts,
	}
	
	// sérialisation de la requête
	url := rva.url_server + "/vote"
	data, _ := json.Marshal(req) // code la requete vote en liste de bit

	// envoi de la requête au url
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// soit le vote a été ajouté soit une erreur est survenu ..

	// traitement de la réponse
	if err != nil {
		
		// A REVOIR [TODO]
		//return "",err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}

	res,err = rva.decodeResponse(resp)
	return
}



func (rva *RestVoterAgent) DoRequestResult(ballotID string) (res rad_t.Response, err error) {
	// creation de requete de resultat
	req := rad_t.RequestVote{
		BallotID: "scrutin1",
	}

	// sérialisation de la requête
	url := rva.url_server + "/result"
	data, _ := json.Marshal(req) // code la requete vote en liste de bit

	// envoi de la requête au url
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}
	res,err = rva.decodeResponse(resp)
	return
}

// TO DO : à vérifier si on mets les ballotID
func (rva *RestVoterAgent) Start(ballotID string) {
	log.Printf("démarrage de %s", rva.agt.ID)
	resp, _ := rva.doRequestVoter(ballotID)
	
	if resp.Status == http.StatusOK{
		log.Print("Vote enregistré !")
	}
}
