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
	url  string //localhost:8080
	opts []int
}

func NewRestVoterAgent(id string, n string, p []coms.Alternative, u string, op []int) *RestVoterAgent {
	ag := rad_t.NewAgent(id, n, p)
	return &RestVoterAgent{ag, u, op}
}

// traduire le résultat en chaine de caractère
func (rva *RestVoterAgent) treatResponseVote(r *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var resp string
	json.Unmarshal(buf.Bytes(), &resp) //decode la reponse et le met dans resp

	return resp
}

/*
renvoie la réponse du serveur ou une erreur
*/
func (rva *RestVoterAgent) doRequestVoter(ballotID string) (res string, err error) {
	// creation de requete de vote
	req := agt.RequestVote{
		AgentID:     string(rva.agt.ID),
		BallotID:    ballotID,
		Preferences: rva.agt.Prefs,
		Options:     rva.opts,
	}

	// sérialisation de la requête
	url := rva.url + "/vote"
	data, _ := json.Marshal(req) // code la requete vote en liste de bit

	// envoi de la requête au url
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// soit le vote a été ajouté soit une erreur est survenu ..

	// traitement de la réponse
	if err != nil {
		// A REVOIR [TODO]
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}
	res = rva.treatResponseVote(resp)

	return
}

// TO DO : à vérifier si on mets les ballotID
func (rva *RestVoterAgent) Start(ballotID string) {
	log.Printf("démarrage de %s", rva.agt.ID)
	_, err := rva.doRequestVoter(ballotID)

	if err != nil {
		log.Fatal(rva.agt.ID, "error:", err.Error())
	} else {
		log.Printf("[%s] Reponse de server au client : Vote enregistree ! \n", rva.agt.ID)
	}
}
