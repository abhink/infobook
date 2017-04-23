package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/sessions"

	"golang.org/x/net/context"
)

const (
	pattern     = "templates/*.html"
	sessionName = "infobook-session"
)

var (
	sessionStore sessions.Store

	tmpl = template.Must(template.ParseGlob(pattern))

	ErrUnauthorized = errors.New("unauthorized")
)

func init() {
	ek := []byte(os.Getenv("SESSION_ENCRYPTION_KEY"))
	ak := []byte(os.Getenv("SESSION_AUTHENTICATION_KEY"))
	if len(ek) == 0 || len(ak) == 0 {
		log.Fatal("Session parameters not set.")
	}

	sessionStore = sessions.NewCookieStore(ak, ek)
}

func checkAccess(ctx context.Context, w http.ResponseWriter, r *http.Request, skipXSRF bool) bool {
	var wg sync.WaitGroup
	var sessionOK, tokenOK bool

	wg.Add(2)

	go func() {
		defer wg.Done()

		err := checkSession(ctx, r)
		if err == nil {
			sessionOK = true
		}
	}()

	go func() {
		defer wg.Done()

		if skipXSRF {
			tokenOK = true
			return
		}

		u, t := r.Form.Get("email"), r.Form.Get("token")
		tokenOK = checkXSRF(u, t)
	}()

	wg.Wait()

	log.Print("CHECKS PASSED: ", sessionOK && tokenOK)
	return sessionOK && tokenOK
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.Handle("/authorise", entrypoint(authHandler))
	http.Handle("/oauthorise", entrypoint(oauthHandler))
	http.Handle("/create/", entrypoint(createHandler))

	http.Handle("/profile/", &endpoint{
		Handler: profileHandler,
		Opts:    handlerOpts{SkipXSRF: true},
	})
	http.Handle("/update/", &endpoint{Handler: updateHandler})

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":"+port, nil)
}

func SetErrorResponse(w http.ResponseWriter, errMsg string, status int) {
	w.WriteHeader(status)

	b, err := json.Marshal(&response{
		ErrorMsg: errMsg,
	})
	if err != nil {
		log.Print("cannot write error response: ", err)
	}
	w.Write(b)
}
