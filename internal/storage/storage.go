package storage

import (
	"context"
	"errors"

	"github.com/okhomin/security/internal/models/acl"

	"github.com/okhomin/security/internal/models/group"

	"github.com/okhomin/security/internal/models/file"

	"github.com/okhomin/security/internal/models/user"
)

var (
	ErrUserNotExist      = errors.New("user doesn't exist")
	ErrUserAlreadyExist  = errors.New("user already exist")
	ErrFileNotExist      = errors.New("file doesn't exist")
	ErrFileAlreadyExist  = errors.New("file already exist")
	ErrGroupNotExist     = errors.New("group doesn't exist")
	ErrGroupAlreadyExist = errors.New("group already exist")
	ErrAclNotExist       = errors.New("acl doesn't exist")
)

type UserWriter interface {
	CreateUser(ctx context.Context, user user.User) (*user.User, error)
}

type UserReader interface {
	User(ctx context.Context, login string) (*user.User, error)
}

type FileWriter interface {
	CreateFile(ctx context.Context, userID string, file file.File) (*file.File, error)
	UpdateFile(ctx context.Context, file file.File) (*file.File, error)
	DeleteFile(ctx context.Context, id string) (*file.File, error)
}

type FileReader interface {
	File(ctx context.Context, id string) (*file.File, error)
	InfosFiles(ctx context.Context) ([]*file.File, error)
	GroupPermissionsFile(ctx context.Context, id, userID string) (bool, bool, error)
	AclPermissionsFile(ctx context.Context, id, userID string) (bool, bool, error)
}

type GroupWriter interface {
	CreateGroup(ctx context.Context, group group.Group) (*group.Group, error)
	DeleteGroup(ctx context.Context, id string) (*group.Group, error)
	UpdateGroup(ctx context.Context, group group.Group) (*group.Group, error)
}

type GroupReader interface {
	ListGroups(ctx context.Context) ([]*group.Group, error)
	Group(ctx context.Context, id string) (*group.Group, error)
}

type AclWriter interface {
	CreateAcl(ctx context.Context, acl acl.Acl) (*acl.Acl, error)
	UpdateAcl(ctx context.Context, acl acl.Acl) (*acl.Acl, error)
	DeleteAcl(ctx context.Context, id string) (*acl.Acl, error)
}

type AclReader interface {
	ListAcls(ctx context.Context) ([]*acl.Acl, error)
	Acl(ctx context.Context, id string) (*acl.Acl, error)
}
