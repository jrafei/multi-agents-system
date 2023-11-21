# Compte-Rendu de TP : AI30

**Contribution** : 
* Jana Rafei
* Julien Pillis

---

Vous trouverez dans ce rapport les différents éléments d'élaboration de notre programme pour le TP sur les méthodes de vote. Celui-ci est découpé en trois parties : la première dédiée l'architecture générale, la deuxième décrit le serveur, les agents (users) ainsi que les ballots de la librairie **agt**  et la dernière partie décrit les méthodes de vote implémentées au sein de la librairie **comsoc**


:::info
Pour visionner correctement ce rapport (images), cliquez **[ici](https://md.picasoft.net/s/yQV4P76CX#)** !
:::
## Architecture générale : 

Pour bien comprendre l'architecture REST de notre programme, nous vous proposons de nous appuyez sur un schéma que nous avons réalisé :

![](https://md.picasoft.net/uploads/upload_9dc7acd06d212a448df0b1077c6a4649.png)

Cette architecture se décompose au sein du paquet ***agt*** : 
- le ballot dans ***ballotagent/ballot.go***
- le serveur dans ***restserveragent/server.go***
- l'utilisateur/agent dans **restvoteragent/agent.go***
- les requêtes propres à l'architecture ***request/request.go***

### Côté client

Relativement simple, nous pouvoir remarqué qu'un agent peut faire une demande de création de ballot ***/new_ballot***, une demande de vote ***/vote*** ou bien une demande de résultat d'un scrutin ***/result***. Bien évidemment, il ne peut pas effectuer les 3 requêtes en même temps (d'où le XOR). Cette requête prend la forme d'une requête http.Request au sein de laquelle les informations utiles sont transportées par une sous-requête spécifique ***RequestBallot*** (si création d'un ballot) ou ***RequestVote*** (pour tout autre action). Cette requête est sérialisée sous le format JSON puis en octets pour son transport.

*Remarque* : L'envoie d'une telle requête est faisable par un agent spécialement implémenté dans notre module, ou bien par un utilisateur disposant des URLs du serveur.


### Côté serveur

Une fois la requête récéptionnée par le serveur, un multiplexeur détermine le méthode à exectuer en fonction de la demande du client : 
````go= 
mux := http.NewServeMux()
mux.HandleFunc("/new_ballot", rsa.init_ballot)
mux.HandleFunc("/vote", rsa.ballotHandler("vote"))
mux.HandleFunc("/result", rsa.ballotHandler("result"))
````

*Note* : rsa = agent serveur

#### Création d'un ballot
Si la requête est un demande de création de vote, on appelle la méthode ***init_ballot***. 

Cette méthode vérifie la conformité des paramètres de la requête et, en si pas de non-conformité, crée un goroutine de la méthode Start() d'un ballot. Le ballot sera donc disponible par le serveur en tout temps, et en parallèle d'un autre traitement du serveur et même d'autres ballots. 
Pour pouvoir communiquer avec le ballot, le serveur sauvegarde le ballot-id qu'il lui a attribué et le channel permettant de communiquer directement avec celui-ci.

```go=
ballots map[string]chan request.RequestVoteBallot
```

Sur ce channel, transitent des requêtes internes ***http.Response***. Pour la création d'un ballot, le serveur retourne simplement au client le ballot ID du nouveau ballot.

```go=
// Requête de réponse générale (transfert via requête http)
type Response struct {
	Ballot_id string `json:"ballot-id,omitempty"`
	Winner    int    `json:"winner,omitempty"`
	Ranking   []int  `json:"ranking,omitempty"`
	Info      string `json:"info_serveur,omitempty"`
}
```

*Remarque* : L'attribut **Info** permet de retourner des informations complémentaires au client, en particulier en cas d'erreur (provenance de l'erreur).

#### Demande de vote et de résultat

Si l'utilisateur souhaite envoyer son vote, la demande passe par un handler (***ballotHandler()***). En effet, le serveur traite les demandes de vote et de résultat de la même façon. Celui-ci vérifie la conformité du ballot_id et de la requête (POST). Si tout est conforme, il transmet la requête au ballot concerné vie une requête ***RequestVoteBallot***, en la transitant par le channel associé au ballot. 

```go=
// Requête pour la prise en compte d'un vote, et le résultat 
// d'un scrutin (transfert via requête http)
type RequestVote struct {
	AgentID     string               `json:"agent-id,omitempty"`
	BallotID    string               `json:"ballot-id"`
	Preferences []comsoc.Alternative `json:"prefs,omitempty"`
	Options     []int                `json:"options,omitempty"`
}

// Requête échangée entre le ballot et le serveur (requete interne)
type RequestVoteBallot struct {
	*RequestVote        //renseigné par le serveur
	Action       string //renseigné par le serveur, ex : vote, result...
	StatusCode   int    //renseigné par le ballot
	Msg          string //renseigné par le ballot
	Winner       int    //renseigné par le ballot
	Ranking      []int  //renseigné par le ballot
}
```

*Remarque* : 
    - RequestVoteBallot hérite/dérive (par un pointeur) afin de conserver les informations transmises par le client
    - Action permet d'informer au ballot l'action demandée par l'utilisateur (vote ou résultat)
    - StatusCode et Msg permettent la propagation des erreurs et des informations de traitement
    - Winner et Ranking retourne les informations liées au résultat d'un scrutin

##### Côté ballot

Si la requête réceptionnée par la méthode Start() du ballot est de type vote, on lance la méthode ***vote()***. Si elle est de type result, on lance la méthode ***result()***. Sinon, une erreur est générée. Chaque méthode vérifie la conformité des paramètres donc elle a bessoin.

```go=
func (rsa *RestBallotAgent) Start() {
	for {
		var resp request.RequestVoteBallot
		req := <-rsa.server_chan
		// Selection de l'action à effectuer
		switch req.Action {
		case "vote":
			resp = rsa.vote(req)
		case "result":
			resp = rsa.result()
		default:
			resp.StatusCode = http.StatusBadRequest
			resp.Msg = "Action inconnue."
		}
		// Transmission de la réponse au serveur
		rsa.server_chan <- resp
	}
}
```

Une fois le bon traitement effectué, le ballot retourne sa réponse (RequestVoteBallot mise à jour) au serveur par le channel.

#### Transmission de la réponse

Lorsque le serveur reçoit la réponse du ballot au sein du handler,
```go=
// Attente de la response du ballot
resp = <-ballot_chan
```

il convertit la réponse RequestVoteBallot du ballot en objet Response, puis la sérialise afin de transmettre une ***http.Response*** au client.


---
## Méthodes de vote (paquet comsoc): 

Vous trouverez dans cette partie les différentes méthodes de vote implémentées : 
- **Majority**
- **Borda**
- **Approval**
- **STV**
- **Copeland**
- **Condorcet**
- **Young/Condorcet**
- **Kramer-Simpson**
- **Dodgson**
- **Kemeny/Young**
    
Ces méthodes sont fonctionnelles et peuvent être utilisées lors de la création d'un ballot.

Les méthodes Majority, Borda, Copeland et Kramer-Simpson sont utilisables à travers les **SWFFactory** et **SCFFactory**. 

Cependant, les méthodes Approval et STV doivent être utilisées par les méthodes **SWFFactoryOptions** et **SCFFactoryOptions** puisque des informations supplémentaires sont nécessaires pour déterminer un vainqueur (respectivement les seuils et un tiebreak)

Enfin, les méthodes de Condorcet, Young/Condorcet, Dodgson et Kemeny-Young **sont utilisables directement** puisqu'elles ne retournent qu'un seul gagant ou aucun. En effet, ces méthodes sont définies pour donner un gagnant, s'il existe, et non un classement ou un ensemble de gagnant.

*Remarque* : le fichier ***comsoc/useful.go*** contient des fonctions qui ne sont pas propres aux méthodes, mais utilisées au sein du paquet.

## Utilisation du programme (paquet cmd)

Le paquet ***cmd*** permet de tester notre architecture. 3 scripts sont proposés : 

- ***launch-all-rest-agents/launch-all-rest-agents.go*** : ce script permet de tester les agents sur la méthode de vote Majority
- ***launch-all-ballots/launch-all-ballots.go*** : ce script permet de tester toutes les méthodes de vote. 
- ***launch-rsagt/launch-rsagt.go*** : ce script démarre le serveur pour que l'on puisse l'utiliser sur un navigateur web (en localhost:8080 dans notre cas).