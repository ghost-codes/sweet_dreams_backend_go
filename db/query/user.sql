-- name: Createuser :one
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
RETURNING *;


-- name: UpdateUser :one
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
RETURNING *;


-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 OR email =$1 LIMIT 1;