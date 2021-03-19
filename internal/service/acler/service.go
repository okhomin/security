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
	return s.reader.ListAcls(ctx)
}

func (s *Service) Create(ctx context.Context, acl acl.Acl) (*acl.Acl, error) {
	return s.writer.CreateAcl(ctx, acl)
}

func (s *Service) Update(ctx context.Context, acl acl.Acl) (*acl.Acl, error) {
	return s.writer.UpdateAcl(ctx, acl)
}

func (s *Service) Read(ctx context.Context, id string) (*acl.Acl, error) {
	return s.reader.Acl(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) (*acl.Acl, error) {
	return s.writer.DeleteAcl(ctx, id)
}
