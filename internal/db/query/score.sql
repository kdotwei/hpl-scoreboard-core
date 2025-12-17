-- name: CreateScore :one
INSERT INTO scores (
  user_id,
  gflops,
  problem_size_n,
  block_size_nb,
  linux_username,
  n,
  nb,
  p,
  q,
  execution_time,
  submitted_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: ListTopScores :many
SELECT * FROM scores
ORDER BY gflops DESC
LIMIT $1;