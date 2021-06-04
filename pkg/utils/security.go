package utils

import "golang.org/x/crypto/bcrypt"

func GenerateFromPassword(password string) string {
	hp, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return ""
	}
	return string(hp)
}

func MatchPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
