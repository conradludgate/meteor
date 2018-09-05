package main

import (
	"log"
	"net/http"
	"time"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Username   string
	RemoteAddr string
}

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("account"); err == nil {
		acc := Account{}
		err = s.Decode("account", cookie.Value, &acc)
		if err == nil && acc.RemoteAddr == r.Header.Get("X-Real-IP") {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	t := r.FormValue("type")
	e := r.FormValue("e") // email
	p := r.FormValue("p") // password
	q := r.FormValue("q") // password confirm

	if t == "login" {

		var hash []byte
		if err := select_hash.QueryRow(e).Scan(&hash); err != nil {
			if e == "" {
				tmpls.ExecuteTemplate(w, "login", loginData{
					e,
					[]alert{
						alert{"Email or password is incorrect", "amber"},
					},
					0,
				})
				return
			}
		}

		if err := bcrypt.CompareHashAndPassword(hash, []byte(p)); err == nil {
			// user logged in correctly
			encoded, err := s.Encode("account", Account{e, r.Header.Get("X-Real-IP")})
			if err == nil {
				cookie := &http.Cookie{
					Name:  "account",
					Value: encoded,

					Expires: time.Now().Add(time.Minute * 60),
					MaxAge:  60 * 60,

					Domain: "mesa.cwp.io",
					Path:   "/",
				}
				http.SetCookie(w, cookie)
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		tmpls.ExecuteTemplate(w, "login", loginData{
			e,
			[]alert{
				alert{"Email or password is incorrect", "amber"},
			},
			0,
		})
		return

	} else if t == "create" {
		if e == "" {
			tmpls.ExecuteTemplate(w, "login", loginData{
				e,
				[]alert{
					alert{"Users must provide a valid email address", "amber"},
				},
				1,
			})
			return
		}

		if p == "" {
			tmpls.ExecuteTemplate(w, "login", loginData{
				e,
				[]alert{
					alert{"Users must provide a password", "amber"},
				},
				1,
			})
			return
		}

		if err := checkmail.ValidateFormat(e); err != nil {
			log.Println("Format error:", err.Error())
			tmpls.ExecuteTemplate(w, "login", loginData{
				e,
				[]alert{
					alert{"Users must provide a valid email address", "amber"},
				},
				1,
			})
			return
		}

		if q != p {
			tmpls.ExecuteTemplate(w, "login", loginData{
				e,
				[]alert{
					alert{"Passwords did not match", "amber"},
				},
				1,
			})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(p), 13)
		if err != nil {
			tmpls.ExecuteTemplate(w, "login", loginData{
				e,
				[]alert{
					alert{"Please enter a different password", "amber"},
				},
				1,
			})
			return
		}

		_, err = insert_acc.Exec(e, hash)
		if err != nil {
			tmpls.ExecuteTemplate(w, "login", loginData{
				e,
				[]alert{
					alert{"Account already exists", "amber"},
				},
				1,
			})
			return
		}

		encoded, err := s.Encode("account", Account{e, r.Header.Get("X-Real-IP")})
		if err == nil {
			cookie := &http.Cookie{
				Name:  "account",
				Value: encoded,

				Expires: time.Now().Add(time.Minute * 60),
				MaxAge:  60 * 60,

				Domain: "mesa.cwp.io",
				Path:   "/",
			}
			http.SetCookie(w, cookie)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	tmpls.ExecuteTemplate(w, "login", loginData{})
}
