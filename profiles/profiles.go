package profiles

import (
	"database/sql"
	"log"

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
	var err error
	db, err = sql.Open("mysql", "abhinav:test@/infobook")
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
