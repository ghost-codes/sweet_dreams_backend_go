// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: user.sql

package db

import (
	"context"
)

const createuser = `-- name: Createuser :one
INSERT INTO users (
  username,
  first_name,
  last_name,
  email,
  hashed_password,
  avatar_url,
  contact,
  security_key
) VALUES (
  $1, $2,$3,$4,$5,$6,$7,$8
)
RETURNING id, username, first_name, last_name, email, hashed_password, avatar_url, contact, security_key, password_changed_at, verified_at, created_at
`

type CreateuserParams struct {
	Username       string  `json:"username"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Email          string  `json:"email"`
	HashedPassword string  `json:"hashed_password"`
	AvatarUrl      *string `json:"avatar_url"`
	Contact        *string `json:"contact"`
	SecurityKey    string  `json:"security_key"`
}

func (q *Queries) Createuser(ctx context.Context, arg CreateuserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createuser,
		arg.Username,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.HashedPassword,
		arg.AvatarUrl,
		arg.Contact,
		arg.SecurityKey,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.HashedPassword,
		&i.AvatarUrl,
		&i.Contact,
		&i.SecurityKey,
		&i.PasswordChangedAt,
		&i.VerifiedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, first_name, last_name, email, hashed_password, avatar_url, contact, security_key, password_changed_at, verified_at, created_at FROM users
WHERE username = $1 OR email =$1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.HashedPassword,
		&i.AvatarUrl,
		&i.Contact,
		&i.SecurityKey,
		&i.PasswordChangedAt,
		&i.VerifiedAt,
		&i.CreatedAt,
	)
	return i, err
}