package utils

import "errors"

var ErrPassowordsDontMatch = errors.New("the passwords don't match")

func SamePasswordVerification(password string, confirmPassword string) error {
	if password == confirmPassword {
		return nil
	} else {
		return ErrPassowordsDontMatch
	}
}
