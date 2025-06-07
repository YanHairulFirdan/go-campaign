-- name: CreateDonatur :one
INSERT INTO donaturs (name, email, user_id, campaign_id)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetPaginatedDonaturs :many
SELECT id, name, email FROM donaturs
WHERE donaturs.campaign_id IN (
        SELECT id FROM campaigns WHERE slug  = $1 AND deleted_at IS NULL 
    ) AND 
    id IN (
        SELECT id FROM payments WHERE status = 5 AND donatur_id = donaturs.id
    )
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCampaignTotalPaidDonaturs :one
SELECT COUNT(*) AS total FROM donaturs
WHERE donaturs.campaign_id IN (
        SELECT id FROM campaigns WHERE slug  = $1 AND deleted_at IS NULL 
    ) AND 
    id IN (
        SELECT id FROM payments WHERE status = 5 AND donatur_id = donaturs.id
    );
