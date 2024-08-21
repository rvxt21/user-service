package enteties

type User struct {
	ID       int
	Name     string
	LastName string
	Role     Role
	Email    string
	Password string
	Login    string
	Status   string
}

type Role string

var (
	Admin Role = "admin"
	Buyer Role = "buyer"
)

func New(password string, email string, login string) *User {
	return &User{
		Password: password,
		Email:    email,
		Login:    login,
	}
}
