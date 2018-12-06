package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var logfile *os.File
var conns []*websocket.Conn

var history []string

// If any errors, just ignore...
func Log(msg ...string) {
	if logfile == nil {
		logfile, _ = os.OpenFile("mesa.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		scan := bufio.NewScanner(logfile)
		for scan.Scan() {
			history = append(history, scan.Text())
		}
	}

	s := strings.Join(msg, " ")

	log.Println(s)

	s = time.Now().Format("2006/01/02 15:04:05") + " " + s + "\n"

	io.WriteString(logfile, s)
	history = append(history, s)

	remove := []*websocket.Conn{}

	for _, conn := range conns {
		err := conn.WriteJSON(WSMessage{LOG, []string{s}})
		if err != nil {
			remove = append(remove, conn)
		}
	}

	for _, conn := range remove {
		RemoveConn(conn)
	}
}

func RegisterConn(conn *websocket.Conn) error {
	if logfile == nil {
		logfile, _ = os.OpenFile("mesa.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		scan := bufio.NewScanner(logfile)
		for scan.Scan() {
			history = append(history, scan.Text())
		}
	}

	trunc := len(history) - 200
	if trunc < 0 {
		trunc = 0	
	}
	
	err := conn.WriteJSON(WSMessage{LOG, history[trunc:]})
	if err != nil {
		return err
	}

	err = conn.WriteJSON(WSMessage{PROC, []string{
		strconv.Itoa(processed),
		strconv.Itoa(toprocess),
	}})
	if err != nil {
		return err
	}

	err = conn.WriteJSON(WSMessage{USER, sessions})
	if err != nil {
		return err
	}

	conns = append(conns, conn)
	return nil
}

func RemoveConn(conn *websocket.Conn) {
	for i, v := range conns {
		if conn == v {
			conns = append(conns[:i], conns[i+1:]...)
			break
		}
	}
}

func LOGClose() {
	if err := logfile.Close(); err != nil {
		log.Println("Error closing Logfile:", err.Error())
	}
	for _, conn := range conns {
		if err := conn.Close(); err != nil {
			log.Println("Error closing websocket:", err.Error())
		}
	}
}
