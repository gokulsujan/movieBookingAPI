package controller

import (
	"golang.org/x/crypto/bcrypt"
)

func PassToHash(pass string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return []byte(""), err
	}
	return hash, nil
}

func HashToPass(hashed string, pass string) bool {
	hashedPass := []byte(hashed)
	err := bcrypt.CompareHashAndPassword(hashedPass, []byte(pass))
	if err != nil {
		// Passwords do not match
		return false
	}
	return true
}
