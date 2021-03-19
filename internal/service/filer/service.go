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
	return s.reader.InfosFiles(ctx)
}

func (s *Service) Create(ctx context.Context, file file.File) (*file.File, error) {
	return s.writer.CreateFile(ctx, file)
}

func (s *Service) Update(ctx context.Context, userID string, file file.File) (*file.File, error) {
	_, writeGroup, err := s.reader.GroupPermissionsFile(ctx, file.ID, userID)
	if err != nil {
		return nil, err
	}

	_, writeAcl, err := s.reader.AclPermissionsFile(ctx, file.ID, userID)
	if err != nil {
		return nil, err
	}

	if !writeAcl && !writeGroup {
		return nil, ErrPermissionDenied
	}
	return s.writer.UpdateFile(ctx, file)
}

func (s *Service) Read(ctx context.Context, userID, id string) (*file.File, error) {
	readGroup, _, err := s.reader.GroupPermissionsFile(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	readAcl, _, err := s.reader.AclPermissionsFile(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if !readAcl && !readGroup {
		return nil, ErrPermissionDenied
	}
	return s.reader.File(ctx, id)
}

func (s *Service) Delete(ctx context.Context, userID, id string) (*file.File, error) {
	_, writeGroup, err := s.reader.GroupPermissionsFile(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	_, writeAcl, err := s.reader.AclPermissionsFile(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if !writeGroup && !writeAcl {
		return nil, ErrPermissionDenied
	}
	return s.writer.DeleteFile(ctx, id)
}
