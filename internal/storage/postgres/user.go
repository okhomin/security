package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/okhomin/security/internal/models/user"
	"github.com/okhomin/security/internal/storage"
)

const getUserQuery = `
SELECT id, login, password FROM users WHERE login = $1;
`

func (p *Postgres) User(ctx context.Context, login string) (*user.User, error) {
	result := new(user.User)
	if err := p.db.QueryRow(ctx, getUserQuery, login).Scan(&result.ID, &result.Login, &result.Password); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrUserNotExist
		}
		return nil, err
	}

	return result, nil
}

const createUserQuery = `
INSERT INTO users (login, password) VALUES ($1, $2)
ON CONFLICT DO NOTHING RETURNING id, login, password;
`

func (p *Postgres) CreateUser(ctx context.Context, u user.User) (*user.User, error) {
	result := new(user.User)
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.QueryRow(ctx, createUserQuery, u.Login, u.Password).Scan(&result.ID, &result.Login, &result.Password); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		if err == pgx.ErrNoRows {
			return nil, storage.ErrUserAlreadyExist
		}
		return nil, err
	}
	if _, err := tx.Exec(ctx, createGroupQuery, u.Login, true, true, []string{result.ID}); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		if err == pgx.ErrNoRows {
			return nil, storage.ErrGroupAlreadyExist
		}
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return result, nil
}
