package grouper

import (
	"context"

	"github.com/okhomin/security/internal/models/group"
)

type Grouper interface {
	List(ctx context.Context) ([]*group.Group, error)
	CreateGroup(ctx context.Context, group group.Group) (*group.Group, error)
	UpdateGroup(ctx context.Context, group group.Group) (*group.Group, error)
	Read(ctx context.Context, id string) (*group.Group, error)
	DeleteGroup(ctx context.Context, id string) (*group.Group, error)
}
