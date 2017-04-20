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
	err := db.QueryRow("SELECT * FROM profiles WHERE email = ?", userId).
		Scan(&r.Email, &r.FullName, &r.Adddress, &r.Telephone)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print("error fetching row: ", err)
			return nil, err
		}
		log.Print("No rows found.")
		return nil, nil
	}
	return &r
}
