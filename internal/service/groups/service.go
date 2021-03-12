package groups

import (
	"github.com/okhomin/security/internal/models/group"
	"github.com/okhomin/security/internal/storage"
)

type Service struct {
	writer storage.GroupWriter
	reader storage.GroupReader
}

func (s *Service) Create(name string, read, write bool, users []string) (*group.Group, error) {
	return nil, nil
}
