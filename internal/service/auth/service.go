package auth

import (
	"context"
	"errors"

	"github.com/okhomin/security/internal/hash"
	"github.com/okhomin/security/internal/storage"

	"github.com/okhomin/security/internal/models"
)

var (
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrAlreadyExist           = errors.New("user is already exist")
)

type Service struct {
	writer storage.Writer
	reader storage.Reader
	hasher hash.Hasher
}

func New(hasher hash.Hasher, writer storage.Writer, reader storage.Reader) *Service {
	return &Service{
		writer: writer,
		reader: reader,
		hasher: hasher,
	}
}

func (s *Service) Login(ctx context.Context, password, login []byte) (*models.User, error) {
	user, err := s.reader.User(ctx, login)
	if err == storage.ErrNotExist {
		return nil, ErrInvalidLoginOrPassword
	}
	if err != nil {
		return nil, err
	}

	valid, err := s.hasher.Compare([]byte(user.Password), password)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, ErrInvalidLoginOrPassword
	}

	return user, nil
}

func (s *Service) Signup(ctx context.Context, password, login []byte) (*models.User, error) {
	hashedPassword, err := s.hasher.Generate(password)
	if err != nil {
		return nil, err
	}

	newUser := models.NewUser(string(login), string(hashedPassword))
	user, err := s.writer.AddUser(ctx, *newUser)
	if err == storage.ErrAlreadyExist {
		return nil, ErrAlreadyExist
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
