package main

import (
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	RemoteAddr string    `json:"-"`
	Expires    time.Time `json:"-"`
	Online     bool      `json:"online"`
}

var sessions map[string]Session

func TimeoutTicker(done chan bool) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			update := false
			for k, v := range sessions {
				if t.After(v.Expires) {
					sessions[k] = Session{
						"",
						time.Unix(0, 0),
						false,
					}
					update = true
				}
			}
			if update {
				for _, conn := range conns {
					conn.WriteJSON(WSMessage{USER, sessions})
				}
			}
		}
	}
}

const SessionAge int = 60 * 60 // 1 hour

func LoginTemplate(w io.Writer, e string, t int, msg string, level string) {
	tmpls.ExecuteTemplate(w, "login", loginData{
		// submit,
		e,
		[]alert{
			alert{msg, level},
		},
		t,
	})
}

func CloseSession(w http.ResponseWriter, email string) {
	cookie := &http.Cookie{
		Name: "account",

		Expires: time.Unix(0, 0),
		MaxAge:  -1,

		Domain: "mesa.cwp.io",
		Path:   "/",
	}
	http.SetCookie(w, cookie)

	if _, ok := sessions[email]; ok {
		sessions[email] = Session{
			"",
			time.Unix(0, 0),
			false,
		}

		for _, conn := range conns {
			conn.WriteJSON(WSMessage{USER, sessions})
		}
	}
}

func MakeSession(w http.ResponseWriter, r *http.Request, email string, seconds int) {
	encoded, err := s.Encode("account", email)
	if err == nil {
		cookie := &http.Cookie{
			Name:  "account",
			Value: encoded,

			Expires: time.Now().Add(time.Duration(seconds) * time.Second),
			MaxAge:  seconds,

			Domain: "mesa.cwp.io",
			Path:   "/",
		}
		http.SetCookie(w, cookie)

		sessions[email] = Session{
			r.Header.Get("X-Real-IP"),
			time.Now().Add(time.Duration(seconds) * time.Second),
			true,
		}

		for _, conn := range conns {
			conn.WriteJSON(WSMessage{USER, sessions})
		}
	}
}

func CheckSession(w http.ResponseWriter, r *http.Request) (t bool, email string) {
	if cookie, err := r.Cookie("account"); err == nil {
		err = s.Decode("account", cookie.Value, &email)
		if err == nil {
			s := sessions[email]
			if s.Online &&
				s.RemoteAddr == r.Header.Get("X-Real-IP") &&
				s.Expires.After(time.Now()) {

				MakeSession(w, r, email, SessionAge)

				return true, email
			}
		}
	}

	CloseSession(w, email)

	return false, email
}

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	if c, _ := CheckSession(w, r); c {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	t := r.FormValue("type")
	e := r.FormValue("e") // email
	p := r.FormValue("p") // password
	q := r.FormValue("q") // password confirm
	// submit := r.FormValue("s") // To Submit

	if t == "login" {

		var hash []byte
		if err := select_hash.QueryRow(e).Scan(&hash); err != nil {
			if e == "" {
				LoginTemplate(w, e, 0, "Email or password is incorrect", "amber")
				return
			}
		}

		if err := bcrypt.CompareHashAndPassword(hash, []byte(p)); err == nil {
			// user logged in correctly
			MakeSession(w, r, e, SessionAge)

			Log("User logged in:", e)

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		LoginTemplate(w, e, 0, "Email or password is incorrect", "amber")
		return

	} else if t == "create" {
		if e == "" {
			LoginTemplate(w, e, 1, "Users must provide a valid email address", "amber")
			return
		}

		if p == "" {
			LoginTemplate(w, e, 1, "Users must provide a password", "amber")
			return
		}

		// var email string
		// if err := select_admin.QueryRow(e).Scan(&email); err != nil || email != e {

		if _, ok := sessions[e]; e != "admin" && !ok {
			LoginTemplate(w, e, 1, "Email address not authorised by admin", "amber")
			return
		}

		if q != p {
			LoginTemplate(w, e, 1, "Passwords did not match", "amber")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(p), 13)
		if err != nil {
			LoginTemplate(w, e, 1, "Please enter a different password", "amber")
			return
		}

		_, err = insert_acc.Exec(e, hash)
		if err != nil {
			LoginTemplate(w, e, 1, "Account already exists", "amber")
			return
		}

		MakeSession(w, r, e, SessionAge)

		Log("Account Created:", e)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// if submit != "" {
	// 	LoginTemplate(w, e, 0, "Session timed out. Please login again", "amber")
	// } else {
	// 	tmpls.ExecuteTemplate(w, "login", loginData{})
	// }

	tmpls.ExecuteTemplate(w, "login", loginData{})
}

func LogoutHandle(w http.ResponseWriter, r *http.Request) {
	_, email := CheckSession(w, r)
	CloseSession(w, email)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
