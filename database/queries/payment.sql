-- name: CreatePayment :one
INSERT INTO payments (transaction_id, donatur_id, donation_id, campaign_id, amount, link, note, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPaymentById :one
SELECT * FROM payments WHERE id = $1;

-- name: GetPaymentByTransactionId :one
SELECT * FROM payments WHERE transaction_id = $1;

-- name: UpdatePaymentFromCallback :one
UPDATE payments
SET 
    status = $2, 
    updated_at = CURRENT_TIMESTAMP,
    vendor = $3,
    method = $4,
    response = $5,
    payment_date = $6
WHERE id = $1
RETURNING *;