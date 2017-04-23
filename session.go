package main

import (
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/context"
)

var xsrf = make(map[string]string)

var m sync.Mutex

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func setSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	userId := r.FormValue("email")

	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Print("error getting session store: ", err)
		return err
	}

	session.Values["userId"] = userId
	if err := session.Save(r, w); err != nil {
		log.Print("error setting session: ", err)
		return err
	}
	log.Print("session set: ", userId)
	return nil
}

func checkSession(ctx context.Context, r *http.Request) error {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Print("error getting session store: ", err)
		return err
	}

	userId, found := session.Values["userId"]
	if !found || userId == "" {
		log.Print("invalid session: ", err)
		return err
	}
	log.Print("SESSION FOUND?? -- ", userId)
	return nil
}

func clearSession(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Print("error getting session store: ", err)
	}

	session.Values["userId"] = ""
	if err := session.Save(r, w); err != nil {
		log.Print("error clearing session: ", err)
	}
}

func setXSRF(userId string) string {
	m.Lock()
	defer m.Unlock()

	token := randStringRunes(8)
	xsrf[userId] = token
	log.Print("XSRF set: ", xsrf)

	return token
}

func checkXSRF(userId, token string) bool {
	m.Lock()
	defer m.Unlock()

	log.Print("XSRF got: ", userId, xsrf)
	return xsrf[userId] == token
}

func clearXSRF(userId string) {
	m.Lock()
	defer m.Unlock()

	delete(xsrf, userId)
}
