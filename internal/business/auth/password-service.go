package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func CheckPassword(loginPass string, registeredPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(registeredPass), []byte(loginPass))
	if err != nil {
		log.Println("the password is not the same | [CheckPassword] ")
		return false
	}
	log.Println("the password is the same | [CheckPassword] ")

	return true
}
