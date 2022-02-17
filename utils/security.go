package utils

import "golang.org/x/crypto/bcrypt"

//HashPassword -> get hashed password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

//VerifyPassword -> check password
func VerifyPassword(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

	return err == nil
}
