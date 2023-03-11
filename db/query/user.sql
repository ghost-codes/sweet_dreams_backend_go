-- name: Createuser :one
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
RETURNING *;


-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 OR email =$1 LIMIT 1;