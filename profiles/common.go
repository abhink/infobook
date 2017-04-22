package profiles

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// secret client id
	clientID     = "526401365242-9ltlqtil676q9jvi815kdkael0nr75fa.apps.googleusercontent.com"
	clientSecret = "HvXRJ4IKE8eHahbryJFG8JRS"
)

// ErrUserExists is an error to denote that the user is already registered.
var ErrUserExists = errors.New("user already exists")

var (
	db                    *sql.DB
	insertCredentialsStmt *sql.Stmt
	insertUserStmt        *sql.Stmt
	updateUserStmt        *sql.Stmt
)

func init() {
	dburl := os.Getenv("DATABASE_URL")

	if dburl == "" {
		log.Fatal("$DATABASE_URL must be set")
	}

	var err error
	db, err = sql.Open("mysql", dburl)
	if err != nil {
		log.Fatal(err)
	}

	insertCredentialsStmt, err = db.Prepare("INSERT INTO credentials(email, pass, type) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	updateUserStmt, err = db.Prepare(`UPDATE profiles SET email=?, name=?, address=?, phone=?
                                WHERE email=?`)
	if err != nil {
		log.Fatal(err)
	}

	insertUserStmt, err = db.Prepare(`INSERT INTO profiles (email, name, address, phone)
                                    VALUES (?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func GetLoginURL() string {
	state := randToken()
	return conf.AuthCodeURL(state)
}

func GetJSON(data json.Marshaler) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Print("error marshalling data: ", err)
		return nil, err
	}
	return b, nil
}

func SetErrorResponse(w http.ResponseWriter, errMsg string, status int) {
	w.WriteHeader(status)
	e := struct {
		ErrorMsg string `json:"errorMsg"`
	}{
		ErrorMsg: errMsg,
	}
	b, err := json.Marshal(e)
	if err != nil {
		log.Print("cannot write error response: ", err)
	}
	w.Write(b)
}
