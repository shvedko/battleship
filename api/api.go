package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pshvedko/battleship/api/websocket"
	"github.com/pshvedko/battleship/battle"
)

type Application struct {
	Logging *log.Logger
	Service battle.Battle
}

type state struct {
	H []byte `json:"H,omitempty"`
}

type point struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

type query struct {
	point
	state
}

type reply struct {
	point
	F int `json:"F"`
	C int `json:"C"`
	state
}

func (a *Application) Begin(w websocket.ResponseWriter, r *websocket.Request) {
	j := json.NewEncoder(w)
	p := a.Service.Begin()
	if p == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		for p.Next() {
			w.WriteHeader(http.StatusOK)
			v := reply{F: p.F(), point: point{X: p.X(), Y: p.Y()}, C: p.C(), state: state{H: p.H()}}
			j.Encode(v)
		}
	}
}

func (a *Application) Click(w websocket.ResponseWriter, r *websocket.Request) {
	var q query
	json.NewDecoder(r.Body).Decode(&q)
	j := json.NewEncoder(w)
	p := a.Service.Click(q.X, q.Y, q.H)
	if p == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		for p.Next() {
			w.WriteHeader(http.StatusOK)
			v := reply{F: p.F(), point: point{X: p.X(), Y: p.Y()}, C: p.C(), state: state{H: p.H()}}
			j.Encode(v)
		}
	}
}

func (a *Application) Reset(w websocket.ResponseWriter, r *websocket.Request) {
	a.Service.Reset()
	w.WriteHeader(http.StatusOK)
}
