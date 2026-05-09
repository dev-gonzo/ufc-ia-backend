package auth

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pepper := os.Getenv("PASSWORD_PEPPER")

	passwordWithPepper := password + pepper

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(passwordWithPepper),
		bcrypt.DefaultCost,
	)

	return string(hash), err
}

func CheckPassword(password string, hash string) bool {
	pepper := os.Getenv("PASSWORD_PEPPER")

	passwordWithPepper := password + pepper

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(passwordWithPepper),
	)

	return err == nil
}
