package main

import (
	"encoding/json"
	"fmt"
	"infobook/profiles"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

type response struct {
	Data     interface{} `json:"data"`
	Token    string      `json:"token"`
	ErrorMsg string      `json:"error"`
}

type handlerOpts struct {
	SkipXSRF bool
}

type endpoint struct {
	Handler func(context.Context, *http.Request) (interface{}, error)
	Opts    handlerOpts
}

// type endpoint func(context.Context, *http.Request) (interface{}, error)

func (e *endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")

	if err := r.ParseForm(); err != nil {
		SetErrorResponse(w, "cannot parse form data", http.StatusInternalServerError)
	}
	access := checkAccess(context.Background(), w, r, e.Opts.SkipXSRF)
	if !access {
		SetErrorResponse(w, "not authorized", http.StatusUnauthorized)
		return
	}

	data, err := e.Handler(context.Background(), r)
	if err != nil {
		s := fmt.Sprintf("error encountered: %s", err.Error())
		SetErrorResponse(w, s, http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(&response{
		Data:  data,
		Token: r.Form.Get("token"),
	})
	if err != nil {
		SetErrorResponse(w, "error fetching user", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

type entrypoint func(context.Context, *http.Request) (interface{}, error)

func (e entrypoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")

	if err := r.ParseForm(); err != nil {
		SetErrorResponse(w, "cannot parse form data", http.StatusInternalServerError)
	}
	if err := setSession(context.Background(), w, r); err != nil {
		SetErrorResponse(w, "internal error", http.StatusInternalServerError)
		return
	}

	token := setXSRF(r.FormValue("email"))

	data, err := e(context.Background(), r)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrUnauthorized {
			status = http.StatusUnauthorized
		}
		SetErrorResponse(w, "internal error", status)
		return
	}

	b, err := json.Marshal(&response{
		Data:  data,
		Token: token,
	})
	if err != nil {
		SetErrorResponse(w, "error fetching user", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func profileHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	email := r.Form.Get("email")
	log.Println("getting for: ", email, r.Form)
	u, err := profiles.GetUserProfile(ctx, email)
	if err != nil {
		panic(err)
	}
	return u, nil
}

func authHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	u, p, ok := r.BasicAuth()
	if !ok || !profiles.CheckAuth(ctx, u, p) {
		return nil, ErrUnauthorized
	}
	return &profiles.User{Email: u}, nil
}

func oauthHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	code := r.Form.Get("code")
	state := r.Form.Get("state")
	log.Println("Got OATUH: ", code, state)

	user, valid := profiles.CheckOAuth(ctx, code)
	if !valid {
		return nil, ErrUnauthorized
	}
	return user, nil
}

func createHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	email := r.Form.Get("email")
	pass := r.Form.Get("pass")
	u, err := profiles.RegisterUser(ctx, email, pass, false)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func updateHandler(ctx context.Context, r *http.Request) (interface{}, error) {
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
	u, err := profiles.UpdateUser(ctx, user)
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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")

	if err := r.ParseForm(); err != nil {
		SetErrorResponse(w, "internal error", http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		clearSession(context.Background(), w, r)
	}()

	go func() {
		defer wg.Done()

		u := r.Form.Get("email")
		clearXSRF(u)
	}()

	wg.Wait()

	http.Redirect(w, r, "/", http.StatusFound)
}
