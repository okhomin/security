package storage

import (
	"context"
	"errors"

	"github.com/okhomin/security/internal/models"
)

var (
	ErrNotExist     = errors.New("user doesn't exist")
	ErrAlreadyExist = errors.New("user already exist")
)

type Writer interface {
	AddUser(ctx context.Context, user models.User) (*models.User, error)
}

type Reader interface {
	User(ctx context.Context, login []byte) (*models.User, error)
}
