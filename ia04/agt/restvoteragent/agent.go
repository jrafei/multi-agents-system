package restvoteragent

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	request "ia04/agt/request"
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
func (*Agent) decodeResponse(r *http.Response) (rep request.Response, err error) {
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
		- 'err' : variable d erreur

======================================
*/
func (agt *Agent) Vote(ballotID string, url_server string) (err error) {
	// creation de requete de vote
	req := request.RequestVote{
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
	log.Println("[", agt.ID, "] Vote request (", ballotID, ") : ", req)

	// soit le vote a été ajouté soit une erreur est survenu ..

	// traitement de la réponse
	if err != nil {
		return err
	}

	// décodage de la réponse
	res, err := agt.decodeResponse(resp)
	log.Println("[", agt.ID, "] Response to vote request (", ballotID, ") : \n		  		Statut : ", res.Status, "\n 		  		Info : ", res.Info)
	if res.Status != http.StatusOK {
		return errors.New(res.Info)
	}
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
	    - 'winner' : gagnant du vote (0 = aucun gagnant)
		- 'ranking' : classement du vote
		- 'err' : variable d erreur

======================================
*/
func (agt *Agent) GetResult(ballotID string, url_server string) (winner comsoc.Alternative, ranking []comsoc.Alternative, err error) {
	winner = 0
	ranking = make([]comsoc.Alternative, 0)

	// creation de requete de resultat
	req := request.RequestVote{
		BallotID: ballotID,
	}

	// sérialisation de la requête
	url := url_server + "/result"
	data, _ := json.Marshal(req) // code la requete vote en liste de bit

	// envoi de la requête au url
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	log.Println("[", agt.ID, "] Result request (", ballotID, ") : ", req)

	// traitement de la réponse
	if err != nil {
		return 0, make([]comsoc.Alternative, 0), err
	}

	res, err := agt.decodeResponse(resp)

	// erreur de decodage
	if err != nil {
		log.Println("[", agt.ID, "] Response to result request (", ballotID, ") : \n			   InvalidUnmarshalError.")
		return
	}

	if res.Status != http.StatusOK {
		log.Println("[", agt.ID, "] Response to result request (", ballotID, ") : \n		  		Statut : ", res.Status, "\n 		  		Info : ", res.Info)
		return winner, ranking, errors.New(res.Info)
	} else {
		log.Println("[", agt.ID, "] Response to result request (", ballotID,
			") : \n		  		Statut : ", res.Status,
			"\n 		  		Info : ", res.Info,
			"\n 		  		Winner : ", res.Winner,
			"\n 		  		Ranking : ", res.Ranking)
		winner = comsoc.Alternative(res.Winner)
		for i, _ := range res.Ranking {
			ranking = append(ranking, comsoc.Alternative(res.Ranking[i]))
		}
		return
	}
}

/*
======================================

	  @brief :
	  'Méthode pour créer un ballot'
	  @params :
	    - 'rule' : méthode de vote
		- 'deadline' : deadline de fin de vote
		- 'voters' : liste des ID des voteurs
		- 'nbAlts' : nombre d'alternatives
		- 'tiebreak' : classement des alternatives pour tiebreak
		- 'url_server' : l'url du serveur qui accueilleura le nouveau ballot
	  @returned :
	    - 'ballot_id' : identifiant du ballot crée
		- 'err' : variable d erreur

======================================
*/
func (agt *Agent) CreateBallot(rule string, deadline string, voters []string, nbAlts int, tiebreak []int, url_server string) (ballot_id string, err error) {
	// créer la requête de création de ballot
	req := request.RequestBallot{
		Rule:     rule,
		Deadline: deadline,
		Voters:   voters,
		Nb_alts:  nbAlts,
		Tiebreak: tiebreak,
	}

	// sérialisation de la requête
	url := url_server + "/new_ballot"
	data, _ := json.Marshal(req) // code la requete vote en liste de bit

	// envoi de la requête au url
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	log.Println("[", agt.ID, "] Ballot creation request : ", req)

	// soit le ballot a été créé, soit une erreur est survenu ..

	// traitement de la réponse
	if err != nil {
		return "", err
	}
	// décodage de la réponse
	res, err := agt.decodeResponse(resp)

	// erreur de décodage
	if err != nil {
		return "",err
	}else if res.Status != http.StatusOK {
		log.Println("[", agt.ID, "] Response to create ballot request : \n		  		    Statut : ", res.Status, "\n 		  		    Info : ", res.Info)
		return "",errors.New(res.Info)
	} else {
		log.Println("[", agt.ID, "] Response to create ballot request :"+
			"\n		  		Ballot id : ", res.Ballot_id,
			"\n 		  		Info : ", res.Info,
			"\n					Statud : ", res.Status)

		return res.Ballot_id, err
	}
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
