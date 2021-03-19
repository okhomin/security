package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/okhomin/security/internal/models/file"
	"github.com/okhomin/security/internal/storage"
)

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

func (p *Postgres) GroupPermissionsFile(ctx context.Context, id, userID string) (bool, bool, error) {
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

func (p *Postgres) AclPermissionsFile(ctx context.Context, id, userID string) (bool, bool, error) {
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

const infosFilesQuery = `
SELECT id, name, groups, acls FROM files;
`

func (p *Postgres) InfosFiles(ctx context.Context) ([]*file.File, error) {
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

func (p *Postgres) CreateFile(ctx context.Context, userID string, f file.File) (*file.File, error) {
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
	if tag, err := p.db.Exec(ctx, updateFileQuery, f.Name, f.Content, f.Groups, f.Acls, f.ID); err != nil {
		return nil, err
	} else if tag.RowsAffected() == 0 {
		return nil, storage.ErrFileNotExist
	}

	return &f, nil
}

const deleteFileQuery = `
DELETE FROM files WHERE id = $1 RETURNING id, name, content, groups, acls;
`

func (p *Postgres) DeleteFile(ctx context.Context, id string) (*file.File, error) {
	result := new(file.File)
	if err := p.db.QueryRow(ctx, deleteFileQuery, id).Scan(&result.ID, &result.Name, &result.Content, &result.Groups, &result.Acls); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrFileNotExist
		}
		return nil, err
	}

	return result, nil
}
