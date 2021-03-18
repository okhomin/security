package acler

import (
	"context"

	"github.com/okhomin/security/internal/models/acl"
)

type Acler interface {
	List(ctx context.Context) ([]*acl.Acl, error)
	Create(ctx context.Context, acl acl.Acl) (*acl.Acl, error)
	Update(ctx context.Context, acl acl.Acl) (*acl.Acl, error)
	Read(ctx context.Context, id string) (*acl.Acl, error)
	Delete(ctx context.Context, id string) (*acl.Acl, error)
}
