package postgres

import (
	"context"

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
