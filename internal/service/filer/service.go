package filer

import (
	"context"

	"github.com/okhomin/security/internal/models/file"
	"github.com/okhomin/security/internal/storage"
)

type Service struct {
	writer storage.FileWriter
	reader storage.FileReader
}

func New(writer storage.FileWriter, reader storage.FileReader) *Service {
	return &Service{
		writer: writer,
		reader: reader,
	}
}

func (s *Service) List(ctx context.Context) ([]*file.File, error) {
	panic("implement me")
}

func (s *Service) Create(ctx context.Context, file file.File) (*file.File, error) {
	panic("implement me")
}

func (s *Service) Update(ctx context.Context, file file.File) (*file.File, error) {
	panic("implement me")
}

func (s *Service) Read(ctx context.Context, id string) (*file.File, error) {
	panic("implement me")
}

func (s *Service) Delete(ctx context.Context, id string) (*file.File, error) {
	panic("implement me")
}
