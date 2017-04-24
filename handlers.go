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
	CheckAccess   bool
	SkipXSRF      bool
	PopulateToken bool
	RefreshToken  bool
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
		SetErrorResponse(w, "access denied", http.StatusUnauthorized)
		return
	}

	data, err := e.Handler(context.Background(), r)
	if err != nil {
		s := fmt.Sprintf("error encountered: %s", err.Error())
		SetErrorResponse(w, s, http.StatusInternalServerError)
		return
	}

	token := r.Form.Get("token")
	if e.Opts.PopulateToken {
		email := r.Form.Get("userid")
		if email != "" {
			token = getXSRF(email)
		}
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

type entrypoint struct {
	Handler func(context.Context, *http.Request) (interface{}, error)
	Opts    handlerOpts
}

func (e *entrypoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")

	if err := r.ParseForm(); err != nil {
		SetErrorResponse(w, "cannot parse form data", http.StatusInternalServerError)
	}

	// Entry point may still require access checks.
	if e.Opts.CheckAccess {
		access := checkAccess(context.Background(), w, r, e.Opts.SkipXSRF)
		if !access {
			SetErrorResponse(w, "access denied", http.StatusUnauthorized)
			return
		}
	}

	if err := setSession(context.Background(), w, r); err != nil {
		s := fmt.Sprintf("error: %s", err)
		SetErrorResponse(w, s, http.StatusInternalServerError)
		return
	}

	token := setXSRF(r.FormValue("email"))

	data, err := e.Handler(context.Background(), r)
	if err != nil {
		if err == ErrUnauthorized {
			SetErrorResponse(w, "access denied", http.StatusUnauthorized)
			return
		}
		s := fmt.Sprintf("error: %s", err)
		SetErrorResponse(w, s, http.StatusInternalServerError)
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
	u, err := profiles.GetUserProfile(ctx, email)
	if err != nil {
		fmt.Errorf("error getting profile for: %s", email)
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

func createHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	email := r.Form.Get("email")
	pass := r.Form.Get("pass")

	// if _, err := mail.ParseAddress(email); err != nil {
	// 	return nil, fmt.Errorf("invalid email address: %v", err)
	// }

	u, err := profiles.RegisterUser(ctx, email, pass, false)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func updateHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	user := &profiles.User{
		r.Form.Get("email"),
		r.Form.Get("name"),
		r.Form.Get("address"),
		r.Form.Get("phone"),
	}
	newID := r.Form.Get("email")
	oldID := r.Form.Get("oldemail")

	if oldID != newID {
		_, err := profiles.ReRegisterUser(ctx, oldID, newID)
		if err != nil {
			log.Printf("error updating credentials: %v", err)
			return nil, err
		}
	}

	u, err := profiles.UpdateUser(ctx, user, oldID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

type templateResponse struct {
	LoginURL string
	Email    string
	Token    string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	param := templateResponse{
		LoginURL: profiles.GetLoginURL(),
	}
	if err := tmpl.ExecuteTemplate(w, "main.html", param); err != nil {
		log.Fatal(err)
	}
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if err := r.ParseForm(); err != nil {
		SetErrorResponse(w, "cannot parse form data", http.StatusInternalServerError)
	}

	code := r.Form.Get("code")

	// TODO: Verify state!!!
	// state := r.Form.Get("state")

	user, err := profiles.RegisterOAuthUser(ctx, code)
	if err != nil {
		SetErrorResponse(w, "user not authorised", http.StatusInternalServerError)
		return
	}

	if err := setSessionValue(ctx, w, r, user.Email); err != nil {
		SetErrorResponse(w, "internal error", http.StatusInternalServerError)
		return
	}

	token := setXSRF(user.Email)

	param := templateResponse{
		Email: user.Email,
		Token: token,
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

func tokenHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	// TODO: Add Header Origin checks.

	return "", nil
}
