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
	agt        *rad_t.Agent
	url_server string //localhost:8080
	opts       []int
}

/*
======================================

	  @brief :
	  'Constructeur de la classe.'
	  @params :
		- 'id' : identifiant unique du voteur
		- 'name' : nom du voteur
		- 'preferences' : préférences du voteur
		- 'url_server' : adresse du serveur
		- 'options' : options supplémentaires
	  @returned :
	    -  Un pointeur sur le voteur créé.

======================================
*/
func NewRestVoterAgent(id string, name string, preferences []coms.Alternative, url_server string, options []int) *RestVoterAgent {
	ag := rad_t.NewAgent(id, name, preferences)
	return &RestVoterAgent{ag, url_server, options}
}

/*
======================================

	  @brief :
	  'Méthode pour décodage une requête réponse.'
	  @params :
		- 'r' : requete http contenant la réponse (binaire)
	  @returned :
	    - 'rep' : structure RequestVote contenant les informations traduites envoyées au client
		- 'err' : variable d erreur

======================================
*/
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

/*
======================================

	  @brief :
	  'Méthode pour effectuer un vote.'
	  @params :
		- 'ballotID' : ID du ballot pour lequel le client vote
	  @returned :
	    - 'res' : réponse retournée par le serveur
		- 'err' : variable d erreur

======================================
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

	res, err = rva.decodeResponse(resp)
	return
}

/*
======================================

	  @brief :
	  'Méthode pour obtenir le résultat d'un ballot.'
	  @params :
		- 'ballotID' : ID du ballot pour lequel le client souhaite le résultat
	  @returned :
	    - 'res' : réponse retournée par le serveur
		- 'err' : variable d erreur

======================================
*/
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
	res, err = rva.decodeResponse(resp)
	return
}

// TO DO : à vérifier si on mets les ballotID
/*
======================================

	@brief :
	'Procédure de mise en fonction d un voteur.'

======================================
*/
func (rva *RestVoterAgent) Start(ballotID string) {
	log.Printf("démarrage de %s", rva.agt.ID)

	resp, _ := rva.doRequestVoter(ballotID)

	if resp.Status == http.StatusOK {
		log.Print("Vote enregistré !")
	}
}
