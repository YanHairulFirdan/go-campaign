-- name: CreateUser :one
INSERT INTO users (name, email, password, created_at, updated_at)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;