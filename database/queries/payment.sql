-- name: GetPaymentById :one
SELECT * FROM payments WHERE id = $1;
