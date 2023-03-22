-- name: CreateApproval :one
INSERT INTO approvals(
  request_id,
  assigned_nurse,
  user_id,
  approved_by,
  status,
  notes
) VALUES( $1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetUserApprovals :many
SELECT * FROM approvals
WHERE user_id=$1
LIMIT $2
OFFSET $3;

-- name: AdminGetUserApprovals :many
SELECT * FROM approvals
LIMIT $1
OFFSET $2;

-- name: DeleteApprovals :exec
DELETE FROM approvals
WHERE id=$1;
