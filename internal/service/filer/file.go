package filer

import (
	"context"

	"github.com/okhomin/security/internal/models/file"
)

type Filer interface {
	List(ctx context.Context) ([]*file.File, error)
	Create(ctx context.Context, file file.File) (*file.File, error)
	Update(ctx context.Context, file file.File) (*file.File, error)
	Read(ctx context.Context, id string) (*file.File, error)
	Delete(ctx context.Context, id string) (*file.File, error)
}
