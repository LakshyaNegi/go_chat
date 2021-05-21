package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type authHandler struct {
	next http.Handler
}

var (
	DBConn *gorm.DB
)

type Users struct {
	gorm.Model
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func initDB() {
	//var err error
	DBConn, err := gorm.Open("sqlite3", "users.db")
	if err != nil {
		panic("Failed to connect database")
	}

	log.Printf("Database connected")

	DBConn.AutoMigrate(&Users{})
	log.Printf("Database migrated")
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
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

	//log.Printf("entered user %s , password %s", username, pass)

	//l := &login{username: username, password: pass}

	// for _, user := range users {
	// 	if user.username == l.username && user.password == l.password {

	// 		http.SetCookie(w, &http.Cookie{
	// 			Name:    "auth",
	// 			Value:   user.username,
	// 			Path:    "/",
	// 			Expires: time.Now().Add(180 * time.Second)})

	// 		w.Header()["Location"] = []string{"/chat"}
	// 		w.WriteHeader(http.StatusTemporaryRedirect)
	// 		return
	// 	}
	// }

	db, err := gorm.Open("sqlite3", "users.db")
	if err != nil {
		panic("Failed to connect database")
	}

	var user Users
	db.Where("Username = ?", username).First(&user)

	//log.Printf("db user %s , password %s", user.Username, user.Password)

	if pass != user.Password {
		log.Printf("Password error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "auth",
		Value:   username,
		Path:    "/",
		Expires: time.Now().Add(180 * time.Second)})

	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusTemporaryRedirect)
	return

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "users.db")
	if err != nil {
		panic("Failed to connect database")
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user := &Users{Username: username, Password: password}
	//log.Printf("login user", user)

	db.Create(user)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}
