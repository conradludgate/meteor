package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	LOG int = iota
	PROC
)

type WSMessage struct {
	Type int      `json:"type"`
	Msg  []string `json:"msg"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func AdminWSHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	RegisterConn(conn)
	conn.SetCloseHandler(func(code int, text string) error {
		RemoveConn(conn)
		return conn.Close()
	})

	for {
		if _, _, err := conn.NextReader(); err != nil {
			RemoveConn(conn)
			conn.Close()
			break
		}
	}
}
