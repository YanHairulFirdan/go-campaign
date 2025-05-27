package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func Password(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) (bool, error) {
	if hashedPassword == "" || password == "" {
		return false, errors.New("hashed password and password cannot be empty")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, errors.New("password does not match")
	}

	return true, nil
}
