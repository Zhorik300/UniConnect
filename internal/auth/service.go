package auth

import (
	"fmt"
	"uniconnect/internal/database"
	"uniconnect/internal/user"

	"golang.org/x/crypto/bcrypt"
)

func Register(name, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`
	_, err = database.DB.Exec(query, name, email, string(hashed))
	return err
}

func Login(email, password string) (*user.User, error) {
	var u user.User
	query := `SELECT id, name, email, password FROM users WHERE email=$1`
	err := database.DB.Get(&u, query, email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &u, nil
}
