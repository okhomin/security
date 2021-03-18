package acler

import (
	"context"

	"github.com/okhomin/security/internal/models/acl"
	"github.com/okhomin/security/internal/storage"
)

type Service struct {
	writer storage.AclWriter
	reader storage.AclReader
}

func New(writer storage.AclWriter, reader storage.AclReader) *Service {
	return &Service{
		writer: writer,
		reader: reader,
	}
}

func (s *Service) List(ctx context.Context) ([]*acl.Acl, error) {
	panic("implement me")
}

func (s *Service) Create(ctx context.Context, acl acl.Acl) (*acl.Acl, error) {
	panic("implement me")
}

func (s *Service) Update(ctx context.Context, acl acl.Acl) (*acl.Acl, error) {
	panic("implement me")
}

func (s *Service) Read(ctx context.Context, id string) (*acl.Acl, error) {
	panic("implement me")
}

func (s *Service) Delete(ctx context.Context, id string) (*acl.Acl, error) {
	panic("implement me")
}
