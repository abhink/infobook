package profiles

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// Google oauth client ID
	clientID     = "526401365242-9ltlqtil676q9jvi815kdkael0nr75fa.apps.googleusercontent.com"
	clientSecret = "HvXRJ4IKE8eHahbryJFG8JRS"
)

// ErrUserExists is an error to denote that the user is already registered.
var ErrUserExists = errors.New("user already exists")

var (
	db                    *sql.DB
	insertCredentialsStmt *sql.Stmt
	updateCredentialsStmt *sql.Stmt
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

	insertCredentialsStmt, err = db.Prepare(`INSERT INTO credentials(userid, password, type)
                                           VALUES(?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}

	updateCredentialsStmt, err = db.Prepare(`UPDATE credentials SET userid=?, password=?, type=?
                                           WHERE userid=?`)
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
