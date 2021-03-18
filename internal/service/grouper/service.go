package grouper

import (
	"context"

	"github.com/okhomin/security/internal/models/group"
	"github.com/okhomin/security/internal/storage"
)

type Service struct {
	writer storage.GroupWriter
	reader storage.GroupReader
}

func New(writer storage.GroupWriter, reader storage.GroupReader) *Service {
	return &Service{
		writer: writer,
		reader: reader,
	}
}

func (s *Service) List(ctx context.Context) ([]*group.Group, error) {
	panic("implement me")
}

func (s *Service) CreateGroup(ctx context.Context, group group.Group) (*group.Group, error) {
	panic("implement me")
}

func (s *Service) UpdateGroup(ctx context.Context, group group.Group) (*group.Group, error) {
	panic("implement me")
}

func (s *Service) Read(ctx context.Context, id string) (*group.Group, error) {
	panic("implement me")
}

func (s *Service) DeleteGroup(ctx context.Context, id string) (*group.Group, error) {
	panic("implement me")
}
