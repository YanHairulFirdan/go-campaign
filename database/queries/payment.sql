-- name: GetPaymentById :one
SELECT * FROM payments WHERE id = $1;

-- name: GetDetailPayment :one
SELECT 
    p.id, p.vendor, p.method, p.link, p.status, p.amount::numeric as amount, p.payment_date,
    d.id as donatur_id, d.name as donatur_name,
    c.id as campaign_id, c.title as campaign_title, c.description as campaign_description, u.name as creator
FROM payments AS p
INNER JOIN
    donaturs as d
    ON p.donatur_id = d.id
INNER JOIN
    campaigns as c
    ON p.campaign_id = c.id
INNER JOIN
    users as u
    ON c.user_id = u.id
WHERE p.transaction_id = $1;


