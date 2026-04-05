package websocket

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type messageWriter struct {
	c      *websocket.Conn
	id     []byte
	mu     chan bool
	buffer bytes.Buffer
	proto  int
	status int
	header http.Header
	dirty  bool
	end    rune
	me     []byte
}

func (w *messageWriter) Header() http.Header {
	return w.header
}

func (w *messageWriter) Write(data []byte) (int, error) {
	w.dirty = true
	return w.buffer.Write(data)
}

func (w *messageWriter) WriteHeader(status int) {
	w.Flush()
	w.dirty = true
	w.status = status
}

func (w *messageWriter) Flush() {
	if !w.dirty || !<-w.mu {
		return
	}
	defer func() { w.mu <- true }()
	m := bytes.Buffer{}
	_, err := m.Write(w.id)
	if err == nil {
		_, err = fmt.Fprintf(&m, "\n%d\n%c\n", w.status, w.end)
		if err == nil {
			_, err = w.buffer.WriteTo(&m)
			if err == nil {
				err = w.c.WriteMessage(w.proto, m.Bytes())
				if err == nil {
					w.dirty = false
					return
				}
			}
		}
	}
	_ = w.c.Close()
}

func (w *messageWriter) End() {
	if w.dirty {
		w.end = '#'
	}
}

func (w *messageWriter) ID() []byte {
	return w.me
}
