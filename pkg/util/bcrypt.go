package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func CheckPassword(password string, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(password),
		[]byte(password),
	)
	return err == nil
}
