package main

import (
	"fmt"
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

func setSessionValue(ctx context.Context, w http.ResponseWriter, r *http.Request, val string) error {
	userID := val

	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Print("error getting session store: ", err)
		return err
	}

	session.Values["userID"] = userID
	if err := session.Save(r, w); err != nil {
		log.Print("error setting session: ", err)
		return err
	}
	log.Print("session set: ", userID)
	return nil
}

func setSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	userID := r.FormValue("email")

	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Print("error getting session store: ", err)
		return err
	}

	session.Values["userID"] = userID
	if err := session.Save(r, w); err != nil {
		log.Print("error setting session: ", err)
		return err
	}
	log.Print("session set: ", userID)
	return nil
}

func checkSession(ctx context.Context, r *http.Request) error {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Print("error getting session store: ", err)
		return err
	}

	userID, found := session.Values["userID"]
	if !found || userID == "" {
		log.Print("invalid session: ", err)
		return fmt.Errorf("invalid session")
	}
	log.Print("SESSION FOUND?? -- ", userID)
	return nil
}

func clearSession(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Print("error getting session store: ", err)
	}

	session.Values["userID"] = ""
	if err := session.Save(r, w); err != nil {
		log.Print("error clearing session: ", err)
	}
}

func setXSRF(userID string) string {
	m.Lock()
	defer m.Unlock()

	token := randStringRunes(8)
	xsrf[userID] = token
	log.Print("XSRF set: ", xsrf)

	return token
}

func getXSRF(userID string) string {
	m.Lock()
	defer m.Unlock()

	return xsrf[userID]
}

func checkXSRF(userID, token string) bool {
	m.Lock()
	defer m.Unlock()

	log.Print("XSRF got: ", userID, xsrf)
	return xsrf[userID] == token
}

func clearXSRF(userID string) {
	m.Lock()
	defer m.Unlock()

	delete(xsrf, userID)
}
