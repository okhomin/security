package models

type User struct {
	ID       string
	Login    string
	Password string
}

func NewUser(login, password string) *User {
	return &User{
		Login:    login,
		Password: password,
	}
}
