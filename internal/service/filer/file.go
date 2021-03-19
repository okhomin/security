package filer

import (
	"context"
	"errors"

	"github.com/okhomin/security/internal/models/file"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrFileAlreadyExist = errors.New("file with such name already exist")
	ErrFileNotExist     = errors.New("file doesn't exist")
)

type Filer interface {
	List(ctx context.Context) ([]*file.File, error)
	Create(ctx context.Context, file file.File) (*file.File, error)
	Update(ctx context.Context, userID string, file file.File) (*file.File, error)
	Read(ctx context.Context, userID, id string) (*file.File, error)
	Delete(ctx context.Context, userID, id string) (*file.File, error)
}
