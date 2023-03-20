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