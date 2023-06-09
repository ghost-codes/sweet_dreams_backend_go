-- name: CreateAdmin :one
INSERT INTO admins(
  username,
  full_name,
  email,
  hashed_password,
  is_super
)VALUES(
    $1,$2,$3,$4,false
) RETURNING *;


-- name: GetAdmin :one
SELECT * FROM admins
WHERE username=$1 OR email=$1;


-- name: GetAdminByID :one
SELECT * FROM admins
WHERE id=$1 ;