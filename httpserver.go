package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"

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

	http.HandleFunc("/", HTMLHandle)
	http.HandleFunc("/login/", LoginHandle)

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	log.Println("Server has loaded successfully")

	log.Fatal(http.ListenAndServe(":6374", nil))
}

func HTMLHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		tmpls.ExecuteTemplate(w, "index", nil)
		return
	}
	// else {
	// 	b, ok := files[r.URL.Path[1:]]
	// 	if !ok {
	// 		w.WriteHeader(http.StatusNotFound)
	// 		return
	// 	}

	// 	w.Write(b)
	// }

	http.ServeFile(w, r, filepath.Join("src", r.URL.Path))
}

type Account struct {
	Username   string
	RemoteAddr string
}

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("account"); err == nil {
		acc := Account{}
		err = s.Decode("account", cookie.Value, &acc)
		if err == nil && acc.RemoteAddr == r.RemoteAddr {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	login := r.FormValue("login")
	create := r.FormValue("create")
	e := r.FormValue("e") // email - no username
	p := r.FormValue("p") // password

	if login == "login" {
		// get hash from database
		// SELECT hash FROM accounts WHERE email=?

		var hash []byte
		if err := select_hash.QueryRow(e).Scan(&hash); err != nil {
			// user is not in DB
			// return login page with email already filled in
			// with hint to create account

			return
		}

		if err := bcrypt.CompareHashAndPassword(hash, []byte(p)); err == nil {
			// user logged in correctly
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	} else if create == "create" {

		q := r.FormValue("q") // password confirm

		if e == "" || q != p {
			return
		}
	}

	tmpls.ExecuteTemplate(w, "login", []alert{})
}
