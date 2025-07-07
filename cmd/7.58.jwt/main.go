package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	db  *sqlx.DB
	err error

	templates = template.Must(template.ParseGlob("templates/*.html"))
)

func main() {
	// this Pings the database trying to connect
	// use sqlx.Open() for sql.Open() semantics
	db, err = sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("GET /{$}", loginPage)
	http.HandleFunc("POST /login/{$}", postLoginPage)

	http.HandleFunc("GET /register/{$}", getRegisterPage)
	http.HandleFunc("POST /register/{$}", postRegisterPage)

	http.HandleFunc("GET /dashboard/{$}", dashboardPage)

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func postLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getRegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "register.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func postRegisterPage(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	pass := r.FormValue("password")
	email := r.FormValue("email")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	_, err = db.NamedExec(`INSERT INTO users (uname, fullname, pass) VALUES (:first,:last,:email)`,
		map[string]interface{}{
			"first": "Bin",
			"last":  "Smuth",
			"email": "bensmith@allblacks.nz",
		})
	if err != nil {
		panic(err)
	}

	if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func dashboardPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "dashboard.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
