-- name: CreateDonation :one
INSERT INTO donations (donatur_id, campaign_id, amount, note, payment_status)
VALUES ($1, $2, $3, $4, $5) RETURNING *;