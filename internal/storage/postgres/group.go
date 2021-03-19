package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/okhomin/security/internal/models/group"
	"github.com/okhomin/security/internal/storage"
)

const getGroupQuery = `
SELECT id, name, read, write, users FROM groups WHERE id = $1;
`

func (p *Postgres) Group(ctx context.Context, id string) (*group.Group, error) {
	result := new(group.Group)
	if err := p.db.QueryRow(ctx, getGroupQuery, id).Scan(&result.ID, &result.Name, &result.Read, &result.Write, &result.Users); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrGroupNotExist
		}
		return nil, err
	}

	return result, nil
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

const deleteGroupQuery = `
DELETE FROM groups WHERE id = $1 RETURNING id, name, read, write, users;
`

func (p *Postgres) DeleteGroup(ctx context.Context, id string) (*group.Group, error) {
	result := new(group.Group)
	if err := p.db.QueryRow(ctx, deleteGroupQuery, id).Scan(&result.ID, &result.Name, &result.Read, &result.Write, &result.Users); err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrGroupNotExist
		}
		return nil, err
	}

	return nil, nil
}

const updateGroupQuery = `
UPDATE groups SET name = $1, read = $2, write = $3, users = $4 WHERE id = $5;
`

func (p *Postgres) UpdateGroup(ctx context.Context, g group.Group) (*group.Group, error) {
	if tag, err := p.db.Exec(ctx, updateFileQuery, g.Name, g.Read, g.Write, g.Users); err != nil {
		return nil, err
	} else if tag.RowsAffected() == 0 {
		return nil, storage.ErrGroupNotExist
	}

	return &g, nil
}

const listGroupsQuery = `
SELECT id, name, read, write, users FROM groups;
`

func (p *Postgres) ListGroups(ctx context.Context) ([]*group.Group, error) {
	result := make([]*group.Group, 0, 10)
	rows, err := p.db.Query(ctx, listGroupsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := new(group.Group)
		if err := rows.Scan(&r.ID, &r.Name, &r.Read, &r.Write, &r.Users); err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}
