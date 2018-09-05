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
	"strconv"

	"github.com/gorilla/securecookie"
)

var (
	hashKey  = securecookie.GenerateRandomKey(64)
	blockKey = securecookie.GenerateRandomKey(32)
	s        = securecookie.New(hashKey, blockKey)
)

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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("Closing...")

			csvw.Flush()
			if err := csvw.Error(); err != nil {
				log.Println("Error flushing CSV:", err.Error())
			}
			if err := file.Close(); err != nil {
				log.Println("Error closing CSV:", err.Error())
			}

			os.Exit(0)
		}
	}()

	http.HandleFunc("/", HTMLHandle)
	http.HandleFunc("/login/", LoginHandle)
	http.HandleFunc("/submit/", SubmitHandle)

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	log.Println("Server has loaded successfully")

	log.Fatal(http.ListenAndServe(":6374", nil))
}

func HTMLHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if cookie, err := r.Cookie("account"); err == nil {
			acc := Account{}
			err = s.Decode("account", cookie.Value, &acc)
			if err != nil || acc.RemoteAddr != r.Header.Get("X-Real-IP") {
				http.Redirect(w, r, "/login/", http.StatusSeeOther)
				return
			}
		} else {
			http.Redirect(w, r, "/login/", http.StatusSeeOther)
			return
		}

		tmpls.ExecuteTemplate(w, "index", getImage())
		return
	}

	http.ServeFile(w, r, filepath.Join("src", r.URL.Path))
}

func SubmitHandle(w http.ResponseWriter, r *http.Request) {
	acc := Account{}
	if cookie, err := r.Cookie("account"); err == nil {
		err = s.Decode("account", cookie.Value, &acc)
		if err != nil || acc.RemoteAddr != r.Header.Get("X-Real-IP") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{
				"error": 1,
				"msg": "Login required"
			}`))
			return
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
			"error": 1,
			"msg": "Login required"
		}`))
		return
	}

	vals := struct {
		Image   string `json:"image"`
		Meteors []struct {
			T int `json:"t"`
			R int `json:"r"`
			B int `json:"b"`
			L int `json:"l"`
		} `json:"meteors"`
	}{}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": 2,
			"msg": "Could not process request: ` + err.Error() + `"
		}`))
		return
	}

	err = json.Unmarshal(b, &vals)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": 2,
			"msg": "Could not process request: ` + err.Error() + `"
		}`))
		return
	}

	if vals.Image == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": 2,
			"msg": "Could not process request: No image provided"
		}`))
		return
	}

	records := [][]string{}

	for i, v := range vals.Meteors {
		record := []string{vals.Image, strconv.Itoa(i),
			strconv.Itoa(v.T),
			strconv.Itoa(v.R),
			strconv.Itoa(v.B),
			strconv.Itoa(v.L),
			acc.Username,
		}
		records = append(records, record)
	}

	err = csvw.WriteAll(records)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
}
