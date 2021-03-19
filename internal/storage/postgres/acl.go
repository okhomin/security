package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/okhomin/security/internal/storage"

	"github.com/okhomin/security/internal/models/acl"
)

const createAclQuery = `
INSERT INTO acls (user_id, read, write) VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING RETURNING id;
`

func (p *Postgres) CreateAcl(ctx context.Context, a acl.Acl) (*acl.Acl, error) {
	if err := p.db.QueryRow(ctx, createAclQuery, a.UserID, a.Read, a.Write).Scan(&a.ID); err != nil {
		return nil, err
	}

	return &a, nil
}

const updateAclQuery = `
UPDATE acls SET user_id = $1, read = $2, write = $3 WHERE id = $4;
`

func (p *Postgres) UpdateAcl(ctx context.Context, acl acl.Acl) (*acl.Acl, error) {
	if tag, err := p.db.Exec(ctx, updateAclQuery, acl.UserID, acl.Read, acl.Write, acl.ID); err != nil {
		return nil, err
	} else if tag.RowsAffected() == 0 {
		return nil, storage.ErrAclNotExist
	}

	return &acl, nil
}

const deleteAclQuery = `
DELETE FROM acls WHERE id = $1 RETURNING id, user_id, read, write;
`

func (p *Postgres) DeleteAcl(ctx context.Context, id string) (*acl.Acl, error) {
	result := new(acl.Acl)
	if err := p.db.QueryRow(ctx, deleteAclQuery, id).Scan(&result.ID, &result.UserID, &result.Read, &result.Write); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrAclNotExist
		}
		return nil, err
	}

	return result, nil
}

const getAclQuery = `
SELECT id, user_id, read, write FROM acls WHERE id = $1;
`

func (p *Postgres) Acl(ctx context.Context, id string) (*acl.Acl, error) {
	result := new(acl.Acl)
	if err := p.db.QueryRow(ctx, getAclQuery, id).Scan(&result.ID, &result.UserID, &result.Read, &result.Write); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrAclNotExist
		}
		return nil, err
	}

	return result, nil
}

const listAclsQuery = `
SELECT id, user_id, read, write FROM acls;
`

func (p *Postgres) ListAcls(ctx context.Context) ([]*acl.Acl, error) {
	result := make([]*acl.Acl, 0, 10)
	rows, err := p.db.Query(ctx, listAclsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := new(acl.Acl)
		if err := rows.Scan(&r.ID, &r.UserID, &r.Read, &r.Write); err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}
