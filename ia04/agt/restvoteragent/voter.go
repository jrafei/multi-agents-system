package restvoteragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	utils "ia04/agt/utils"
	comsoc "ia04/comsoc"
)

type AgentID string

type AgentI interface {
	Equal(ag AgentI) bool
	DeepEqual(ag AgentI) bool
	Clone() AgentI
	String() string
	Prefers(a comsoc.Alternative, b comsoc.Alternative) bool
	Start()
}

type Agent struct {
	ID    AgentID
	Name  string
	Prefs []comsoc.Alternative
	Opts  []int
}

/*
======================================

	  @brief :
	  'Constructeur de la classe.'
	  @params :
		- 'id' : identifiant unique du voteur
		- 'name' : nom du voteur
		- 'preferences' : préférences du voteur
		- 'options' : options supplémentaires
	  @returned :
	    -  Un pointeur sur le voteur créé.

======================================
*/
func NewAgent(id string, name string, preferences []comsoc.Alternative, options []int) *Agent {
	return &Agent{AgentID(id), name, preferences, options}
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
func (*Agent) decodeResponse(r *http.Response) (rep utils.Response, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &rep)
	rep.Status = r.StatusCode
	return
}

/*
======================================

	  @brief :
	  'Méthode pour effectuer un vote.'
	  @params :
		- 'ballotID' : ID du ballot pour lequel le client vote
		- 'url_server' : l'url du serveur accueillant le ballot
	  @returned :
	    - 'res' : réponse retournée par le serveur
		- 'err' : variable d erreur

======================================
*/
func (agt *Agent) Vote(ballotID string, url_server string) (res utils.Response, err error) {
	// creation de requete de vote
	req := utils.RequestVote{
		AgentID:     string(agt.ID),
		BallotID:    ballotID,
		Preferences: agt.Prefs,
		Options:     agt.Opts,
	}

	// sérialisation de la requête
	url := url_server + "/vote"
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

	res, err = agt.decodeResponse(resp)
	return
}

/*
======================================

	  @brief :
	  'Méthode pour obtenir le résultat d'un ballot.'
	  @params :
		- 'ballotID' : ID du ballot pour lequel le client souhaite le résultat
		- 'url_server' : l'url du serveur accueillant le ballot
	  @returned :
	    - 'res' : réponse retournée par le serveur
		- 'err' : variable d erreur

======================================
*/
func (agt *Agent) GetResult(ballotID string, url_server string) (res utils.Response, err error) {
	// creation de requete de resultat
	req := utils.RequestVote{
		BallotID: "scrutin1",
	}

	// sérialisation de la requête
	url := url_server + "/result"
	data, _ := json.Marshal(req) // code la requete vote en liste de bit

	// envoi de la requête au url
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}
	res, err = agt.decodeResponse(resp)
	return
}

/* =====================METHODES SUPPLEMENTAIRES=========================== */

func (a *Agent) Equal(ag Agent) bool {
	return a == &ag
}

func (a *Agent) DeepEqual(ag Agent) bool {
	return a.ID == ag.ID && a.Name == ag.Name && slicesEquality[comsoc.Alternative](a.Prefs, ag.Prefs) && slicesEquality[int](a.Opts, ag.Opts)
}

func slicesEquality[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (a *Agent) Clone() Agent {
	prefs_slc := make([]comsoc.Alternative, len(a.Prefs))
	for i, v := range a.Prefs {
		prefs_slc[i] = v
	}
	opts_slc := make([]int, len(a.Opts))
	for i, v := range a.Opts {
		opts_slc[i] = v
	}
	return Agent{a.ID, a.Name, prefs_slc, opts_slc}
}

func (a *Agent) String() string {
	var infos string
	infos = "--------------------------\n"
	infos += "Agent ID : " + string(a.ID) + "\n"
	infos += "Agent name : " + a.Name + "\n"
	infos += "Agent preferences : \n"
	for i, v := range a.Prefs {
		infos += strconv.Itoa(i) + "." + strconv.Itoa(int(v)) + "\n"
	}
	infos += "Agent options : \n"
	for i, v := range a.Opts {
		infos += strconv.Itoa(i) + "." + strconv.Itoa(v) + "\n"
	}

	infos += "-------------------------"
	return infos
}

func (ag *Agent) Prefers(a comsoc.Alternative, b comsoc.Alternative) bool {
	for _, v := range ag.Prefs {
		if v == a {
			return true
		} else if v == b {
			return false
		}
	}
	return false
}

/* ==================================================================== */
