package service

type storage interface {
	CreateUser(string, string) error
}
type Service struct {
	S storage
}

func (s Service) SignUp(password string, email string) error {
	err := s.S.CreateUser(password, email)
	if err != nil {
		return err
	}
	return nil
}
