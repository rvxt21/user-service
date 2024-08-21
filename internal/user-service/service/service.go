package service

type storage interface {
	CreateUser(string, string, string) error
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
