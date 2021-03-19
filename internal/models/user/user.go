package user

type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func New(login, password string) *User {
	return &User{
		Login:    login,
		Password: password,
	}
}
