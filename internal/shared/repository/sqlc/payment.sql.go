// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: payment.sql

package sqlc

import (
	"context"
)

const getPaymentById = `-- name: GetPaymentById :one
SELECT id, transaction_id, donatur_id, donation_id, campaign_id, vendor, method, amount, link, note, status, response, payment_date, created_at, updated_at FROM payments WHERE id = $1
`

func (q *Queries) GetPaymentById(ctx context.Context, id int32) (Payment, error) {
	row := q.db.QueryRowContext(ctx, getPaymentById, id)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.TransactionID,
		&i.DonaturID,
		&i.DonationID,
		&i.CampaignID,
		&i.Vendor,
		&i.Method,
		&i.Amount,
		&i.Link,
		&i.Note,
		&i.Status,
		&i.Response,
		&i.PaymentDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
