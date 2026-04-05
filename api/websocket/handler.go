package websocket

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocket struct {
	u websocket.Upgrader
	h mux.Router
}

func New() *WebSocket {
	return &WebSocket{}
}

func (s *WebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := s.u.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	z := make(chan bool, 1)
	z <- true
	defer close(z)
	defer func() { <-z }()
	var t int
	var b []byte
	var q *http.Request
	for {
		t, b, err = c.ReadMessage()
		if err != nil {
			return
		}
		switch t {
		case websocket.CloseMessage:
			return
		case websocket.TextMessage:
			p := bytes.Split(b, []byte{'\n'})
			var o io.Reader
			var m []byte
			var u []byte
			switch len(p) {
			case 4:
				o = bytes.NewReader(p[3])
				fallthrough
			case 3:
				u = bytes.TrimSpace(p[2])
				m = bytes.TrimSpace(p[1])
				switch string(m) {
				case http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete:
				default:
					return
				}
			default:
				return
			}
			q, err = http.NewRequestWithContext(r.Context(), string(m), string(u), o)
			if err != nil {
				return
			}
			go func(w *messageWriter, r *http.Request) {
				s.h.ServeHTTP(w, r)
				w.End()
				w.Flush()
			}(&messageWriter{
				c:      c,
				id:     p[0],
				mu:     z,
				proto:  t,
				buffer: bytes.Buffer{},
				status: http.StatusOK,
				header: http.Header{},
				dirty:  false,
				end:    '+',
			}, q)
		}
	}
}

func (s *WebSocket) UpgradeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connection := r.Header.Get("Connection")
		upgrade := r.Header.Get("Upgrade")
		if strings.EqualFold(connection, "Upgrade") && strings.EqualFold(upgrade, "Websocket") {
			s.ServeHTTP(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
}

type Request struct {
	*http.Request
}

type HandlerFunc func(w ResponseWriter, r *Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w.(ResponseWriter), &Request{Request: r})
}

func (s *WebSocket) HandleFunc(path string, f func(w ResponseWriter, r *Request)) {
	s.h.Handle(path, HandlerFunc(f))
}
