package user

type User struct {
	ID       string
	Login    string
	Password string
}

func New(login, password string) *User {
	return &User{
		Login:    login,
		Password: password,
	}
}
