package auther

import (
	"context"
	"errors"

	"github.com/okhomin/security/internal/models/user"
)

var (
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrAlreadyExist           = errors.New("user already exist")
)

type Auther interface {
	Login(ctx context.Context, password, login string) (*user.User, error)
	Signup(ctx context.Context, password, login string) (*user.User, error)
}
