package domain

type User struct {
	Id       string
	Username string
	Email    string
	Password string
}

type UsersService interface {
	SignUp(User) (string, error)
	SignIn(User) (string, error)
	Authenticate(token string) (string, error)
	UserExists(userId string) (bool, error)
}

type UsersRepository interface {
	Create(User) (string, error)
	GetByEmail(email string) (User, error)
	CheckExistence(userId string) (bool, error)
}
