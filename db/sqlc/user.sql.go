// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: user.sql

package db

import (
	"context"
	"time"
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
  security_key,
  twitter_social,
  google_social,
  apple_social 
) VALUES (
  $1, $2,$3,$4,$5,$6,$7,$8,$9,$10,$11
)
RETURNING id, username, first_name, last_name, email, hashed_password, avatar_url, contact, security_key, password_changed_at, verified_at, created_at, twitter_social, google_social, apple_social
`

type CreateuserParams struct {
	Username       string  `json:"username"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Email          string  `json:"email"`
	HashedPassword *string `json:"hashed_password"`
	AvatarUrl      *string `json:"avatar_url"`
	Contact        *string `json:"contact"`
	SecurityKey    string  `json:"security_key"`
	TwitterSocial  bool    `json:"twitter_social"`
	GoogleSocial   bool    `json:"google_social"`
	AppleSocial    bool    `json:"apple_social"`
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
		arg.TwitterSocial,
		arg.GoogleSocial,
		arg.AppleSocial,
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
		&i.TwitterSocial,
		&i.GoogleSocial,
		&i.AppleSocial,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, first_name, last_name, email, hashed_password, avatar_url, contact, security_key, password_changed_at, verified_at, created_at, twitter_social, google_social, apple_social FROM users
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
		&i.TwitterSocial,
		&i.GoogleSocial,
		&i.AppleSocial,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users 
 SET username = $2,
  first_name = $3,
  last_name = $4,
  email = $5,
  hashed_password = $6,
  avatar_url = $7,
  contact = $8,
  security_key = $9,
  password_changed_at = $10,
  verified_at = $11,
  created_at = $12,
   twitter_social=$13,
  google_social=$14,
  apple_social=$15 
WHERE id = $1
RETURNING id, username, first_name, last_name, email, hashed_password, avatar_url, contact, security_key, password_changed_at, verified_at, created_at, twitter_social, google_social, apple_social
`

type UpdateUserParams struct {
	ID                int64      `json:"id"`
	Username          string     `json:"username"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Email             string     `json:"email"`
	HashedPassword    *string    `json:"hashed_password"`
	AvatarUrl         *string    `json:"avatar_url"`
	Contact           *string    `json:"contact"`
	SecurityKey       string     `json:"security_key"`
	PasswordChangedAt time.Time  `json:"password_changed_at"`
	VerifiedAt        *time.Time `json:"verified_at"`
	CreatedAt         time.Time  `json:"created_at"`
	TwitterSocial     bool       `json:"twitter_social"`
	GoogleSocial      bool       `json:"google_social"`
	AppleSocial       bool       `json:"apple_social"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.Username,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.HashedPassword,
		arg.AvatarUrl,
		arg.Contact,
		arg.SecurityKey,
		arg.PasswordChangedAt,
		arg.VerifiedAt,
		arg.CreatedAt,
		arg.TwitterSocial,
		arg.GoogleSocial,
		arg.AppleSocial,
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
		&i.TwitterSocial,
		&i.GoogleSocial,
		&i.AppleSocial,
	)
	return i, err
}
