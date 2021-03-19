package grouper

import (
	"context"

	"github.com/okhomin/security/internal/models/group"
)

type Grouper interface {
	List(ctx context.Context) ([]*group.Group, error)
	Create(ctx context.Context, group group.Group) (*group.Group, error)
	Update(ctx context.Context, group group.Group) (*group.Group, error)
	Read(ctx context.Context, id string) (*group.Group, error)
	Delete(ctx context.Context, id string) (*group.Group, error)
}
