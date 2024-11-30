// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: menu.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createMenu = `-- name: CreateMenu :one
INSERT INTO menus (
  name, start_date, end_date
) VALUES (
  $1, $2, $3
)
RETURNING id, name, start_date, end_date
`

type CreateMenuParams struct {
	Name      string
	StartDate pgtype.Date
	EndDate   pgtype.Date
}

func (q *Queries) CreateMenu(ctx context.Context, arg CreateMenuParams) (Menu, error) {
	row := q.db.QueryRow(ctx, createMenu, arg.Name, arg.StartDate, arg.EndDate)
	var i Menu
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.StartDate,
		&i.EndDate,
	)
	return i, err
}

const deleteMenu = `-- name: DeleteMenu :exec
DELETE FROM menus
WHERE id = $1
`

func (q *Queries) DeleteMenu(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteMenu, id)
	return err
}

const getMenu = `-- name: GetMenu :one
SELECT id, name, start_date, end_date FROM menus
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetMenu(ctx context.Context, id int64) (Menu, error) {
	row := q.db.QueryRow(ctx, getMenu, id)
	var i Menu
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.StartDate,
		&i.EndDate,
	)
	return i, err
}

const listMenus = `-- name: ListMenus :many
SELECT id, name, start_date, end_date FROM menus
ORDER BY start_date
`

func (q *Queries) ListMenus(ctx context.Context) ([]Menu, error) {
	rows, err := q.db.Query(ctx, listMenus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Menu
	for rows.Next() {
		var i Menu
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.StartDate,
			&i.EndDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMenu = `-- name: UpdateMenu :one
UPDATE menus
set name = $2, start_date = $3, end_date = $4
  WHERE id = $1
RETURNING id, name, start_date, end_date
`

type UpdateMenuParams struct {
	ID        int64
	Name      string
	StartDate pgtype.Date
	EndDate   pgtype.Date
}

func (q *Queries) UpdateMenu(ctx context.Context, arg UpdateMenuParams) (Menu, error) {
	row := q.db.QueryRow(ctx, updateMenu,
		arg.ID,
		arg.Name,
		arg.StartDate,
		arg.EndDate,
	)
	var i Menu
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.StartDate,
		&i.EndDate,
	)
	return i, err
}