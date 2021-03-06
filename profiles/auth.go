package profiles

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var env = ""

var conf = &oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	RedirectURL:  "http://127.0.0.1:5000/oauthorise",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}

func init() {
	env = os.Getenv("ENV")
	if env == "" {
		env = "prod"
	}

	log.Print("Env set: ", env)
	if env == "prod" {
		conf.RedirectURL = "https://stark-fjord-40589.herokuapp.com/oauthorise"
	}
}

func getCredentials(ctx context.Context, userId string) (string, string, string, error) {
	var u, p, t string
	if err := db.QueryRow("SELECT * FROM credentials WHERE userid = ?", userId).
		Scan(&u, &p, &t); err != nil {
		if err != sql.ErrNoRows {
			log.Print("error fetching row: ", err)
			return "", "", "", err
		}
		log.Print("No rows found.")
		return "", "", "", nil
	}
	return u, p, t, nil
}

func CheckAuth(ctx context.Context, userId, pass string) bool {
	_, p, t, err := getCredentials(ctx, userId)
	if err != nil {
		log.Print("error getting credentials: ", err)
		return false
	}
	if t == "GOOGLE" {
		log.Print("credentials type oauth: ", err)
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(p), []byte(pass)); err != nil {
		return false
	}
	return true
}

func CheckOAuth(ctx context.Context, code string) ([]byte, bool) {
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Print("OAuth code exchange failed: ", err)
		return nil, false
	}
	// Construct the client.
	client := conf.Client(oauth2.NoContext, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Print("OAuth client fetch failed: ", err)
		return nil, false
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(data))
	return data, true
}
