package main

import (
	"log"
	"net/http"
	"time"
)

type authHandler struct {
	next http.Handler
}

type login struct {
	username string
	password string
}

var users = []login{
	login{username: "lak", password: "lak"},
	login{username: "abc", password: "abc"},
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error)
	} else {
		h.next.ServeHTTP(w, r)
	}

}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("auth"); err == http.ErrNoCookie {
	} else if err != nil {
		panic(err.Error())
	} else {
		c.Name = "Deleted"
		c.Value = "Unuse"
		c.Expires = time.Now()
	}

	if r.Method != http.MethodPost {
		log.Printf("Login method was not post")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	pass := r.FormValue("password")

	l := &login{username: username, password: pass}

	for _, user := range users {
		if user.username == l.username && user.password == l.password {

			http.SetCookie(w, &http.Cookie{
				Name:    "auth",
				Value:   user.username,
				Path:    "/",
				Expires: time.Now().Add(180 * time.Second)})

			w.Header()["Location"] = []string{"/chat"}
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
	}

	log.Printf("User not found")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
