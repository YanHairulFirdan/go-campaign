-- name: CreateDonatur :one
INSERT INTO donaturs (name, email, user_id, campaign_id)
VALUES ($1, $2, $3, $4) RETURNING *;