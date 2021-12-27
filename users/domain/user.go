package domain

type User struct {
	Username string
	Email    string
	Password string
}

type UsersService interface {
	SignUp(User) (string, error)
	SignIn(User) (string, error)
}

type UsersRepository interface {
	Create(User) (string, error)
	GetByEmail(User) (string, string, error)
}
