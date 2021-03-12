package auth

import (
	"context"
	"errors"

	"github.com/okhomin/security/internal/models/user"

	"github.com/okhomin/security/internal/hash"
	"github.com/okhomin/security/internal/storage"
)

var (
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrAlreadyExist           = errors.New("user already exist")
)

type Service struct {
	writer storage.UserWriter
	reader storage.UserReader
	hasher hash.Hasher
}

func New(hasher hash.Hasher, writer storage.UserWriter, reader storage.UserReader) *Service {
	return &Service{
		writer: writer,
		reader: reader,
		hasher: hasher,
	}
}

func (s *Service) Login(ctx context.Context, password, login []byte) (*user.User, error) {
	user, err := s.reader.User(ctx, login)
	if err == storage.ErrUserNotExist {
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

func (s *Service) Signup(ctx context.Context, password, login []byte) (*user.User, error) {
	hashedPassword, err := s.hasher.Generate(password)
	if err != nil {
		return nil, err
	}

	newUser := user.New(string(login), string(hashedPassword))
	user, err := s.writer.CreateUser(ctx, *newUser)
	if err == storage.ErrUserAlreadyExist {
		return nil, ErrAlreadyExist
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
