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
	return s.reader.ListGroups(ctx)
}

func (s *Service) Create(ctx context.Context, group group.Group) (*group.Group, error) {
	return s.writer.CreateGroup(ctx, group)
}

func (s *Service) Update(ctx context.Context, group group.Group) (*group.Group, error) {
	return s.writer.UpdateGroup(ctx, group)
}

func (s *Service) Read(ctx context.Context, id string) (*group.Group, error) {
	return s.reader.Group(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) (*group.Group, error) {
	return s.writer.DeleteGroup(ctx, id)
}
