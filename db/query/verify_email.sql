-- name: CreateVerifyEmail :one
INSERT INTO verify_emails(
  username,
  email ,
  secret_key
) VALUES (
  $1, $2,$3
)
RETURNING *;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET is_used=true
WHERE id=$1
RETURNING *;

-- name: GetVerifyEmail :one
SELECT * FROM verify_emails
WHERE id=$1;