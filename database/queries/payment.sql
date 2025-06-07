-- name: CreatePayment :one
INSERT INTO payments (transaction_id, donatur_id, donation_id, campaign_id, amount, link, note, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;