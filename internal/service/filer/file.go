package filer

import (
	"context"
	"errors"

	"github.com/okhomin/security/internal/models/file"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
)

type Filer interface {
	List(ctx context.Context) ([]*file.File, error)
	Create(ctx context.Context, userID string, file file.File) (*file.File, error)
	Update(ctx context.Context, userID string, file file.File) (*file.File, error)
	Read(ctx context.Context, userID, id string) (*file.File, error)
	Delete(ctx context.Context, userID, id string) (*file.File, error)
}
