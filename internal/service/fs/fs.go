package fs

import (
	"github.com/okhomin/security/internal/storage"
)

type Service struct {
	writer storage.FileWriter
	reader storage.FileReader
}
