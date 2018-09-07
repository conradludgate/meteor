package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	LOG int = iota
	PROC
	USER
)

type WSMessage struct {
	Type int         `json:"type"`
	Msg  interface{} `json:"msg"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func AdminWSHandle(w http.ResponseWriter, r *http.Request) {
	c, email := CheckSession(w, r)

	if !c || email != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
	}

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
		_, r, err := conn.NextReader()
		if err != nil {
			RemoveConn(conn)
			conn.Close()
			break
		}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			continue
		}
		email := string(b)
		insert_admin.Exec(email)
	}
}
