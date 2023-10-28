package agt

import (
	"strconv"
	coms "ia04/comsoc"
)


type AgentID string

type AgentI interface {
	Equal(ag AgentI) bool
	DeepEqual(ag AgentI) bool
	Clone() AgentI
	String() string
	Prefers(a coms.Alternative, b coms.Alternative) bool
	Start()
}

type Agent struct {
	ID    AgentID
	Name  string
	Prefs []coms.Alternative
}


func NewAgent(id string, n string, prf []coms.Alternative) *Agent{
	return &Agent{AgentID(id),n,prf}
}
func (a *Agent) Equal(ag Agent) bool {
	return a == &ag
}

func (a *Agent) DeepEqual(ag Agent) bool {
	return a.ID == ag.ID && a.Name == ag.Name && slicesEquality(a.Prefs, ag.Prefs)
}

func slicesEquality(a, b []coms.Alternative) bool {
	// Vérifie l'égalité de deux slices
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
	prefs_slc := make([]coms.Alternative, len(a.Prefs))
	for i, _ := range a.Prefs {
		prefs_slc[i] = a.Prefs[i]
	}
	return Agent{a.ID, a.Name, prefs_slc}
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
	infos += "-------------------------"
	return infos
}


func (ag *Agent) Prefers(a coms.Alternative, b coms.Alternative) bool {
	for _, v := range ag.Prefs {
		if v == a {
			return true
		} else if v == b {
			return false
		}
	}
	return false // Par défaut, à vérifier
}

func Start(){};
