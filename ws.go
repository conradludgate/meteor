package main

import (
	"encoding/json"
	"net/http"
	"time"

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

type WSRequest struct {
	Type int    `json:"type"`
	Data string `json:"data"`
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

		var wsr WSRequest
		err = json.NewDecoder(r).Decode(&wsr)
		if err != nil {
			continue
		}

		if wsr.Type == 0 {
			_, err := insert_admin.Exec(email)
			if err == nil {
				sessions[email] = Session{
					"",
					time.Unix(0, 0),
					false,
				}
				for _, conn := range conns {
					conn.WriteJSON(WSMessage{USER, sessions})
				}
				Log("Added user", email)
			}
		} else if wsr.Type == 1 {
			_, err := delete_admin.Exec(email, email)
			if err == nil {
				delete(sessions, email)
				for _, conn := range conns {
					conn.WriteJSON(WSMessage{USER, sessions})
				}
				Log("Deleted user", email)
			}
		}
	}
}
