package profiles

import (
	"database/sql"
	"log"
	"os"

	"golang.org/x/net/context"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
	Email     string
	FullName  string
	Adddress  string
	Telephone string
}

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
}

func GetUserProfile(ctx context.Context, userId string) (*User, error) {
	r := User{}
	err := db.QueryRow("SELECT * FROM profiles WHERE email = ?", userId).Scan(&r)
	if err != nil {
		log.Print("error fetching orw: ", err)
		return nil, err
	}
	return &r, nil
}
