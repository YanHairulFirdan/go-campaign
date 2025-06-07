-- name: CreateDonation :one
INSERT INTO donations (donatur_id, campaign_id, amount, note)
VALUES ($1, $2, $3, $4) RETURNING *;