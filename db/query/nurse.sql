-- name: CreateNurse :one
INSERT INTO nurses(
  full_name,
  email,
  contact,
  profile_picture
)
VALUES($1,$2,$3,$4)
RETURNING *;


-- name: GetNurse :one
SELECT * FROM nurses
WHERE id=$1;

-- name: ListNurses :many
SELECT * FROM nurses
ORDER BY full_name
LIMIT $1
OFFSET $2;

-- name: DeleteNurse :exec
DELETE FROM nurses
WHERE id=$1;
