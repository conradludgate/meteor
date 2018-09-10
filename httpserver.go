package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/gorilla/securecookie"
)

var (
	hashKey  = securecookie.GenerateRandomKey(64)
	blockKey = securecookie.GenerateRandomKey(32)
	s        = securecookie.New(hashKey, blockKey)
)

func shutdown() {
	Log("Shutting down server")

	CSVClose()
	SQLClose()
	LOGClose()

	os.Exit(0)
}

func main() {
	if err := openSqlDB("acc.db"); err != nil {
		fmt.Println("Could not open database.\n")
		fmt.Println(err.Error())
		return
	}

	if err := sqlPrepareStmts(); err != nil {
		fmt.Println("Could not load prepared statements.")
		fmt.Println(err.Error())
		return
	}

	if err := loadSrc(); err != nil {
		fmt.Println("Could not load source files.")
		fmt.Println(err.Error())
		return
	}

	if err := openCSV("data.csv"); err != nil {
		fmt.Println("Could not open CSV file.")
		fmt.Println(err.Error())
		return
	}

	done := make(chan bool, 1)
	go TimeoutTicker(done)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			done <- true
			shutdown()
		}
	}()

	http.HandleFunc("/", HTMLHandle)
	http.HandleFunc("/login/", LoginHandle)
	http.HandleFunc("/logout/", LogoutHandle)
	http.HandleFunc("/submit/", SubmitHandle)
	http.HandleFunc("/admin/", AdminHandle)
	http.HandleFunc("/adminws/", AdminWSHandle)

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	Log("Server has loaded successfully")

	log.Fatal(http.ListenAndServe(":6374", nil))
}

func HTMLHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		c, email := CheckSession(w, r)
		if !c {
			http.Redirect(w, r, "/login/", http.StatusSeeOther)
			return
		}

		if email == "admin" {
			http.Redirect(w, r, "/admin/", http.StatusSeeOther)
			return
		}

		tmpls.ExecuteTemplate(w, "index", getImage())
		return
	}

	http.ServeFile(w, r, filepath.Join("src", r.URL.Path))
}

func SubmitHandle(w http.ResponseWriter, r *http.Request) {
	c, email := CheckSession(w, r)
	if !c {
		// w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
			"error": 1,
			"msg": "Login required"
		}`))
		return
	}

	vals := submitData{}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": 2,
			"msg": "Could not process request: ` + err.Error() + `"
		}`))
		return
	}

	err = json.Unmarshal(b, &vals)
	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": 2,
			"msg": "Could not process request: ` + err.Error() + `"
		}`))
		return
	}

	if vals.Image == "" {
		// w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": 2,
			"msg": "Could not process request: No image provided"
		}`))
		return
	}

	err = submit(vals, email)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{
			"error": 3,
			"msg": "Could not save data: ` + err.Error() + `"
		}`))
		return
	}

	w.Write([]byte(`{
		"error": 0,
		"msg": "` + getImage() + `"
	}`))

	Log("Data submitted by", email, "for", vals.Image)
}

func AdminHandle(w http.ResponseWriter, r *http.Request) {
	c, email := CheckSession(w, r)
	if !c {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	if email != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpls.ExecuteTemplate(w, "admin", nil)
	return
}
