package profiles

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func UpdateUser(ctx context.Context, user *User, oldID string) (*User, error) {
	_, err := updateUserStmt.Exec(user.Email, user.FullName, user.Address, user.Phone, oldID)
	if err != nil {
		log.Print("error updating user: ", err)
		return nil, err
	}
	return user, nil
}

func RegisterUser(ctx context.Context, userId, pass string, oauth bool) (*User, error) {
	u, _, _, err := getCredentials(ctx, userId)
	if err != nil {
		log.Print("error getting credentials: ", err)
		return nil, fmt.Errorf("error verifying user")
	}
	if u != "" {
		return nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Print("error creating account: ", err)
		return nil, fmt.Errorf("error creating user account")
	}

	_, err = insertCredentialsStmt.Exec(userId, hash, "LOGIN")
	if err != nil {
		log.Print("error creating user: ", err)
		return nil, fmt.Errorf("error creating user account")
	}

	_, err = insertUserStmt.Exec(userId, "", "", "")
	if err != nil {
		log.Print("error creating user, deregistering: ", err)
		deregisterUser(ctx, userId)
		return nil, fmt.Errorf("error creating user account")
	}
	return &User{Email: userId}, nil
}

func ReRegisterUser(ctx context.Context, userID, newID string) (*User, error) {
	u, _, _, err := getCredentials(ctx, newID)
	if err != nil {
		log.Print("error getting credentials: ", err)
		return nil, fmt.Errorf("error verifying user")
	}
	if u != "" {
		return nil, ErrUserExists
	}

	u, p, t, err := getCredentials(ctx, userID)
	if err != nil {
		log.Print("error getting credentials: ", err)
		return nil, fmt.Errorf("error verifying user")
	}
	if u == "" || t == "GOOGLE" {
		return nil, fmt.Errorf("update not allowed: %v", userID)
	}

	_, err = updateCredentialsStmt.Exec(newID, p, "LOGIN", userID)
	if err != nil {
		log.Print("error updating user: ", err)
		return nil, fmt.Errorf("error updating user account")
	}

	return &User{Email: newID}, nil
}

func RegisterOAuthUser(ctx context.Context, code string) (*User, error) {
	data, ok := CheckOAuth(ctx, code)
	if !ok {
		return nil, fmt.Errorf("error verifying user")
	}

	user := User{}
	err := json.Unmarshal(data, &user)
	if err != nil {
		return nil, fmt.Errorf("invalid user data")
	}

	u, _, t, err := getCredentials(ctx, user.Email)
	if err != nil {
		log.Print("error getting credentials: ", err)
		return nil, fmt.Errorf("error getting credentials")
	}
	if u != "" && t != "GOOGLE" {
		log.Print("credentials type login: ", err)
		return nil, fmt.Errorf("credentials of type login")
	}

	if u == "" {
		_, err = insertCredentialsStmt.Exec(user.Email, "", "GOOGLE")
		if err != nil {
			log.Print("error creating user: ", err)
			return nil, fmt.Errorf("error creating user account")
		}

		_, err = insertUserStmt.Exec(user.Email, user.FullName, "", "")
		if err != nil {
			log.Print("error creating user, deregistering: ", err)
			deregisterUser(ctx, user.Email)
			return nil, fmt.Errorf("error creating user account")
		}
	}
	return &user, nil
}

func deregisterUser(ctx context.Context, userId string) {
	_, err := db.Exec("DELETE FROM credentials WHERE userid = ?", userId)
	if err != nil {
		log.Fatal("Failed to deregister user: ", err)
	}
}
