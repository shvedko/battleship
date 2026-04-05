package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pshvedko/battleship/api"
	"github.com/pshvedko/battleship/api/websocket"
	"github.com/pshvedko/battleship/battle"
)

//go:embed html
var h embed.FS

func main() {
	k := [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1}
	b := battle.New(1, 4, 3, 3, 2, 2, 2, 1, 1, 1, 1)
	err := b.Encryption(k)
	if err != nil {
		log.Fatal(err)
	}
	a := api.Application{
		Battle: b,
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	w := websocket.New()
	w.HandleFunc("/begin", a.Begin)
	w.HandleFunc("/click", a.Click)
	w.HandleFunc("/reset", a.Reset)
	r := mux.NewRouter()
	f, err := fs.Sub(h, "html")
	if err != nil {
		log.Fatal(err)
	}
	r.PathPrefix("/").Handler(http.FileServer(http.FS(f))).Methods(http.MethodGet, http.MethodHead)
	r.Use(a.LoggingMiddleware)
	r.Use(w.UpgradeMiddleware)
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
