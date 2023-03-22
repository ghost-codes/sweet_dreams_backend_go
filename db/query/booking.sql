-- name: CreateBookingRequest :one
INSERT INTO requests(
  user_id,
  type,
  prefered_nurse,
  start_date,
  end_date,
  location
)
VALUES($1,$2,$3,$4,$5,Point($6,$7))
RETURNING *;

-- name: ListUserBookingReqs :many
SELECT * FROM requests
WHERE user_id=$1
LIMIT $2
OFFSET $3;

-- name: GetBookingByID :one
SELECT * FROM requests
WHERE user_id=$1 AND id=$2;

-- name: DeleteBookingByID :exec
DELETE FROM requests
WHERE user_id=$1 AND id=$2;

-- name: GetAllBookingsByAdmin :many
SELECT * FROM requests;

-- name: GetBookingsByAdminByID :one
SELECT * FROM requests
WHERE id=$1;

