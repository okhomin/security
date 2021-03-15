package postgres

import (
	"context"

	"github.com/okhomin/security/internal/models/acl"

	"github.com/okhomin/security/internal/models/group"

	"github.com/okhomin/security/internal/models/file"

	"github.com/okhomin/security/internal/models/user"

	"github.com/okhomin/security/internal/storage"

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

const createGroupQuery = `
INSERT INTO groups (name, read, write, users) VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING RETURNING id, name, read, write, users;
`

func (p *Postgres) CreateGroup(ctx context.Context, g group.Group) (*group.Group, error) {
	result := new(group.Group)
	if err := p.db.QueryRow(ctx, createGroupQuery, g.Name, g.Read, g.Write, g.Users).Scan(&result.ID, &result.Name, &result.Read, &result.Write, &result.Users); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrGroupAlreadyExist
		}
		return nil, err
	}

	return result, nil
}

const getGroupQuery = `
SELECT id, name, read, write, users FROM groups WHERE name = $1;
`

func (p *Postgres) Group(ctx context.Context, name string) (*group.Group, error) {
	result := new(group.Group)
	if err := p.db.QueryRow(ctx, getGroupQuery, name).Scan(&result.ID, &result.Name, &result.Read, &result.Write, &result.Users); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrGroupNotExist
		}
		return nil, err
	}

	return result, nil
}

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

const getFileQuery = `
SELECT id, name, content, groups, acls FROM files WHERE id = $1;
`

func (p *Postgres) File(ctx context.Context, id string) (*file.File, error) {
	result := new(file.File)
	if err := p.db.QueryRow(ctx, getFileQuery, id).Scan(&result.ID, &result.Name, &result.Content, &result.Groups, &result.Acls); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrFileNotExist
		}
		return nil, err
	}

	return result, nil
}

const permissionsGroupFileQuery = `
SELECT g.read, g.write FROM groups AS g WHERE 
EXISTS (SELECT 1 FROM files AS f WHERE g.id = ANY (f.groups) AND f.id = $1)
AND $2 = any (g.users);
`

func (p *Postgres) FileGroupPermissions(ctx context.Context, id, userID string) (bool, bool, error) {
	var read, write bool

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return false, false, err
	}
	row, err := tx.Query(ctx, getFileQuery, id)

	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return false, false, err
		}
		if err == pgx.ErrNoRows {
			return false, false, storage.ErrFileNotExist
		}
		return false, false, err
	}
	row.Close()
	if err := tx.QueryRow(ctx, permissionsGroupFileQuery, id, userID).Scan(&read, &write); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return false, false, err
		}
		if err == pgx.ErrNoRows {
			return false, false, nil
		}
		return false, false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return false, false, err
	}

	return read, write, nil
}

const permissionsAclFileQuery = `
SELECT a.read, a.write FROM acls AS a WHERE
EXISTS (SELECT 1 FROM files AS f WHERE a.id = ANY (f.acls) AND f.id = $1)
AND a.user_id = $2;
`

func (p *Postgres) FileAclPermissions(ctx context.Context, id, userID string) (bool, bool, error) {
	var read, write bool

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return false, false, err
	}
	row, err := tx.Query(ctx, getFileQuery, id)

	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return false, false, err
		}
		if err == pgx.ErrNoRows {
			return false, false, storage.ErrFileNotExist
		}
		return false, false, err
	}
	row.Close()
	if err := tx.QueryRow(ctx, permissionsAclFileQuery, id, userID).Scan(&read, &write); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return false, false, err
		}
		if err == pgx.ErrNoRows {
			return false, false, nil
		}
		return false, false, err
	}
	if err := tx.Rollback(ctx); err != nil {
		return false, false, err
	}

	return read, write, nil
}

const createAclQuery = `
INSERT INTO acls (user_id, read, write) VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING RETURNING id, user_id, read, write;
`

func (p *Postgres) CreateAcl(ctx context.Context, userID string, read, write bool) (*acl.Acl, error) {
	result := new(acl.Acl)
	if err := p.db.QueryRow(ctx, createAclQuery, userID, read, write).Scan(&result.ID, &result.UserID, &result.Read, &result.Write); err != nil {
		return nil, err
	}

	return result, nil
}

const infosFilesQuery = `
SELECT id, name, groups, acls FROM files;
`

func (p *Postgres) FilesInfos(ctx context.Context) ([]*file.File, error) {
	result := make([]*file.File, 0, 10)
	rows, err := p.db.Query(ctx, infosFilesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		f := new(file.File)
		if err := rows.Scan(&f.ID, &f.Name, &f.Groups, &f.Acls); err != nil {
			return nil, err
		}

		result = append(result, f)
	}

	return result, nil
}

const createFileQuery = `
INSERT INTO files (name, content, groups, acls) VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING RETURNING id, name, content, groups, acls;
`

func (p *Postgres) CreateFile(ctx context.Context, f file.File) (*file.File, error) {
	result := new(file.File)
	if err := p.db.QueryRow(ctx, createFileQuery, f.Name, f.Content, f.Groups, f.Acls).Scan(&result.ID, &result.Name, &result.Content, &result.Groups, &result.Acls); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrFileAlreadyExist
		}
		return nil, err
	}

	return result, nil
}

const updateFileQuery = `
UPDATE files SET name = $1, content = $2, groups = $3, acls = $4 WHERE id = $5;
`

func (p *Postgres) UpdateFile(ctx context.Context, f file.File) (*file.File, error) {
	if _, err := p.db.Query(ctx, updateFileQuery, f.Name, f.Content, f.Groups, f.Acls, f.ID); err != nil {
		return nil, err
	}

	return &f, nil
}
