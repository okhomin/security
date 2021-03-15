package storage

import (
	"context"
	"errors"

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
)

type UserWriter interface {
	CreateUser(ctx context.Context, user user.User) (*user.User, error)
}

type UserReader interface {
	User(ctx context.Context, login string) (*user.User, error)
}

type FileWriter interface {
	CreateFile(ctx context.Context, file file.File) (*file.File, error)
}

type FileReader interface {
	File(ctx context.Context, name string) (*file.File, error)
	FilesInfos(ctx context.Context) ([]*file.File, error)
	Permissions(ctx context.Context, name, id string) (bool, bool, error)
}

type GroupWriter interface {
	CreateGroup(ctx context.Context, group group.Group) (*group.Group, error)
}

type GroupReader interface {
	Group(ctx context.Context, name string) (*group.Group, error)
}
