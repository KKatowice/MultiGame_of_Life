// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package sqlc

import (
	"context"
	"database/sql"
)

const allUser = `-- name: AllUser :many
SELECT userid, name FROM users
`

func (q *Queries) AllUser(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, allUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.Userid, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createGame = `-- name: CreateGame :exec
INSERT INTO games (roomId, hei, wid) VALUES (?,?,?)
`

type CreateGameParams struct {
	Roomid int32
	Hei    sql.NullInt32
	Wid    sql.NullInt32
}

func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) error {
	_, err := q.db.ExecContext(ctx, createGame, arg.Roomid, arg.Hei, arg.Wid)
	return err
}

const createLobby = `-- name: CreateLobby :exec
INSERT INTO lobby (userId, roomId, id) VALUES (?,?,?)
`

type CreateLobbyParams struct {
	Userid sql.NullInt32
	Roomid sql.NullInt32
	ID     int32
}

func (q *Queries) CreateLobby(ctx context.Context, arg CreateLobbyParams) error {
	_, err := q.db.ExecContext(ctx, createLobby, arg.Userid, arg.Roomid, arg.ID)
	return err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (userId, name) VALUES (?,?)
`

type CreateUserParams struct {
	Userid int32
	Name   sql.NullString
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser, arg.Userid, arg.Name)
	return err
}

const listUsers = `-- name: ListUsers :one
SELECT userId FROM users WHERE users.name = (?)
`

func (q *Queries) ListUsers(ctx context.Context, name sql.NullString) (int32, error) {
	row := q.db.QueryRowContext(ctx, listUsers, name)
	var userid int32
	err := row.Scan(&userid)
	return userid, err
}
