-- name: CreateScore :one
INSERT INTO scores (
  user_id,
  gflops,
  problem_size_n,
  block_size_nb,
  submitted_at
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListTopScores :many
SELECT * FROM scores
ORDER BY gflops DESC
LIMIT $1;