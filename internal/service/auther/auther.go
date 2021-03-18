package auther

import (
	"context"

	"github.com/okhomin/security/internal/models/user"
)

type Auther interface {
	Login(ctx context.Context, password, login string) (*user.User, error)
	Signup(ctx context.Context, password, login string) (*user.User, error)
}
