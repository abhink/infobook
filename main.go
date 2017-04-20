package main

import (
	"html/template"
	"infobook/profiles"
	"log"
	"net/http"
	"path"

	"golang.org/x/net/context"
)

const (
	pattern = "templates/*.html"
)

var tmpl = template.Must(template.ParseGlob(pattern))

func profileHandler(w http.ResponseWriter, r *http.Request) {
	email := path.Base(r.URL.Path)
	u, err := profiles.GetUserProfile(context.Background(), email)
	if err != nil {
		panic(err)
	}
	param := struct {
		Email string
		User  *profiles.User
	}{
		Email: email,
		User:  u,
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", param); err != nil {
		log.Fatal(err)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	if err := tmpl.ExecuteTemplate(w, "temp.html", r.Form); err != nil {
		log.Fatal(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/authorise", authHandler)
	http.HandleFunc("/profile/", profileHandler)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}
