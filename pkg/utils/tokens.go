package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var MySigningKey = []byte("secret-key-users")

func CreateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenStr, err := token.SignedString(MySigningKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
