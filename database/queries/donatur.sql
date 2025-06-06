-- name: CreateDonatur :one
INSERT INTO donaturs (name, email, user_id, campaign_id, amount)
VALUES ($1, $2, $3, $4, $5) RETURNING *;