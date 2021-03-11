package postgres

import (
	"context"

	"github.com/okhomin/security/internal/storage"

	"github.com/okhomin/security/internal/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
)

type Postgres struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, url string) *Postgres {
	db, err := pgxpool.Connect(ctx, url)
	if err != nil {
		panic(err)
	}

	conn, err := db.Acquire(ctx)
	if err != nil {
		panic(err)
	}
	if err := migrateDatabase(ctx, conn.Conn()); err != nil {
		panic(err)
	}
	return &Postgres{
		db: db,
	}
}

func migrateDatabase(ctx context.Context, conn *pgx.Conn) error {
	migrator, err := migrate.NewMigrator(ctx, conn, "users_version")
	if err != nil {
		return err
	}

	if err := migrator.LoadMigrations("./migrations"); err != nil {
		return err
	}

	if err := migrator.Migrate(ctx); err != nil {
		return err
	}
	return nil
}

const getUserQuery = `
SELECT id, login, password FROM users WHERE login = $1;
`

func (p *Postgres) User(ctx context.Context, login []byte) (*models.User, error) {
	result := new(models.User)
	if err := p.db.QueryRow(ctx, getUserQuery, login).Scan(&result.ID, &result.Login, &result.Password); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrNotExist
		}
		return nil, err
	}

	return result, nil
}

const addUserQuery = `
INSERT INTO users (login, password) VALUES ($1, $2)
ON CONFLICT DO NOTHING RETURNING id, login, password;
`

func (p *Postgres) AddUser(ctx context.Context, user models.User) (*models.User, error) {
	result := new(models.User)

	if err := p.db.QueryRow(ctx, addUserQuery, user.Login, user.Password).Scan(&result.ID, &result.Login, &result.Password); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrAlreadyExist
		}
		return nil, err
	}

	return result, nil
}
