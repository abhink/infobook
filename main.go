package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"infobook/profiles"
)

const (
	pattern     = "templates/*.html"
	sessionName = "infobook-session"
)

var (
	// sessionStore sessions.Store

	tmpl = template.Must(template.ParseGlob(pattern))
)

// func init() {
// 	ek := []byte(os.Getenv("SESSION_ENCRYPTION_KEY"))
// 	ak := []byte(os.Getenv("SESSION_AUTHENTICATION_KEY"))
// 	if len(ex) == 0 || len(ak) == 0 {
// 		log.Fatal("Session parameters not set.")
// 	}

// 	sessionStore = sessions.NewCookieStore(ak, ek)
// }

type endpoint func(http.ResponseWriter, *http.Request) (interface{}, error)

type handler struct {
	f         endpoint
	protected bool
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.protected {
		access := checkAccess(context.Background(), r)
		if !access {
			profiles.SetErrorResponse(w, "not authorized", http.StatusUnauthorized)
			return
		}
	}
	data, err := h.f(w, r)
	if err != nil {
		s := fmt.Sprintf("error encountered: %s", err.Error())
		profiles.SetErrorResponse(w, s, http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(data)
	if err != nil {
		profiles.SetErrorResponse(w, "error fetching user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func profileHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	email := r.Form.Get("id")
	log.Println("getting for: ", email, r.Form)
	u, err := profiles.GetUserProfile(context.Background(), email)
	if err != nil {
		panic(err)
	}
	return u, nil
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()
	if !ok || !profiles.CheckAuth(context.Background(), u, p) {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))
		return
	}

	w.Write([]byte(fmt.Sprintf("Authorised: %s - %s\n", u, p)))
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	code := r.Form.Get("code")
	state := r.Form.Get("state")
	log.Println("Got OATUH: ", code, state)

	user, valid := profiles.CheckOAuth(context.Background(), code)
	if !valid {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))
		return
	}
	w.Write(user)
}

func createHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	email := r.Form.Get("email")
	pass := r.Form.Get("pass")
	u, err := profiles.RegisterUser(context.Background(), email, pass, false)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func updateHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	user := &profiles.User{
		r.Form.Get("email"),
		r.Form.Get("name"),
		r.Form.Get("address"),
		r.Form.Get("phone"),
	}
	log.Print("Updating: ", user)
	u, err := profiles.UpdateUser(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	param := struct {
		LoginURL string `json`
	}{
		LoginURL: profiles.GetLoginURL(),
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", param); err != nil {
		log.Fatal(err)
	}
}

func checkAccess(ctx context.Context, r *http.Request) bool {
	log.Print("Called access check")
	return true
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/authorise", authHandler)
	http.HandleFunc("/oauthorise", oauthHandler)
	// http.HandleFunc("/profile/", profileHandler)
	http.Handle("/profile/", &handler{
		protected: true,
		f:         endpoint(profileHandler),
	})
	http.Handle("/create/", &handler{
		f: endpoint(createHandler),
	})
	http.Handle("/update/", &handler{
		protected: true,
		f:         endpoint(updateHandler),
	})
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":"+port, nil)
}
