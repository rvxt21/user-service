package service

import "user/internal/user-service/handlers"

type storage interface {
	CreateUser(handlers.SignUpRequestBody) error
}
type Service struct {
	S storage
}

func (s Service) SignUp(signUpBody handlers.SignUpRequestBody) error {
	err := s.S.CreateUser(signUpBody)
	if err != nil {
		return err
	}
	return nil
}
