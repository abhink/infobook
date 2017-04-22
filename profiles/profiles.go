package profiles

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type User struct {
	Email    string `json:"email"`
	FullName string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

func GetUserProfile(ctx context.Context, userId string) (*User, error) {
	log.Print("calling with:", userId)
	r := User{}
	err := db.QueryRow("SELECT * FROM profiles WHERE email = ?", userId).
		Scan(&r.Email, &r.FullName, &r.Address, &r.Phone)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print("error fetching row: ", err)
			return nil, err
		}
		log.Print("No rows found.")
		return &User{Email: userId}, nil
	}
	return &r, nil
}

func UpdateUser(ctx context.Context, user *User) (*User, error) {
	_, err := updateUserStmt.Exec(user.Email, user.FullName, user.Address, user.Phone, user.Email)
	if err != nil {
		log.Print("error updating user: ", err)
		return nil, err
	}
	return user, nil
}

func RegisterUser(ctx context.Context, userId, pass string, oauth bool) (*User, error) {
	t := "LOGIN"
	if oauth {
		t = "GOOGLE"
	}
	u, _, _, err := getCredentials(ctx, userId)
	if err != nil {
		log.Print("error getting credentials: ", err)
		return nil, err
	}
	if u != "" {
		return nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Print("error creating account: ", err)
		return nil, err
	}

	_, err = insertCredentialsStmt.Exec(userId, hash, t)
	if err != nil {
		log.Print("error creating user: ", err)
		return nil, err
	}

	_, err = insertUserStmt.Exec(userId, "", "", "")
	if err != nil {
		log.Print("error creating user, deregistering: ", err)
		deregisterUser(ctx, userId)
		return nil, err
	}
	return &User{Email: userId}, nil
}

func deregisterUser(ctx context.Context, userId string) {
	_, err := db.Exec("DELETE FROM credentials WHERE userid = ?", userId)
	if err != nil {
		log.Fatal("Failed to deregister user: ", err)
	}
}
