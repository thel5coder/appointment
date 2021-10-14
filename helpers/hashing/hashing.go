package hashing

import (
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hash), err
}

func CheckHashString(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
