package service

import (
	"user/internal/user-service/enteties"
)

type storage interface {
	CreateUser(string, string, string) error
	GetPasswordByEmail(email string) (string, error)
	GetUserByEmail(email string) (enteties.UserPersonalInfo, error)
}
type Service struct {
	S storage
}

func (s Service) SignUp(password, email, login string) error {
	err := s.S.CreateUser(password, email, login)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) SignIn(email string) (string, error) {
	storedPassword, err := s.S.GetPasswordByEmail(email)
	return storedPassword, err
}

func (s Service) GetPersonalInfo(userEmail string) (enteties.UserPersonalInfo, error) {
	userInfo, err := s.S.GetUserByEmail(userEmail)
	if err != nil {
		return enteties.UserPersonalInfo{}, err
	}
	return userInfo, nil
}
