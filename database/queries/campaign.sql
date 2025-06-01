-- name: GetPaginatedUserCampaign :many
SELECT id, title, 
	   CASE 
		   WHEN current_amount = 0 THEN 0 
		   ELSE target_amount / current_amount 
	   END::DECIMAL(10, 2) AS progress, 
	   start_date, end_date, status,
	   CASE
	   	   	WHEN status = 0 THEN 'Draft'
	   	   	WHEN status = 1 THEN 'Active'
	   	   	WHEN status = 2 THEN 'Completed'
	   	   	WHEN status = 3 THEN 'Cancelled'
	   	   ELSE 'Unknown'
	   END AS status_label
FROM campaigns
WHERE 
	user_id = $1 AND
	deleted_at IS NULL AND
	title ILIKE '%' || sqlc.arg(title)::text || '%' AND
	status = sqlc.arg(status)::integer
ORDER BY start_date DESC
LIMIT $2 OFFSET $3;

-- name: GetUserCampaignById :one
SELECT * FROM campaigns
WHERE id = $1 AND user_id = $2;

-- name: CreateCampaign :one
INSERT INTO campaigns (title, description, slug, user_id, target_amount, start_date, end_date, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UpdateCampaign :one
UPDATE campaigns
SET title = $1, description = $2, slug = $3, target_amount = $4, start_date = $5, end_date = $6, status = $7, updated_at = CURRENT_TIMESTAMP
WHERE id = $8 AND user_id = $9
RETURNING id, title, description, slug, user_id, target_amount, current_amount, start_date, end_date, status, created_at::TIMESTAMP, updated_at::TIMESTAMP;

-- name: SoftDeleteCampaign :one
UPDATE campaigns
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: GetCampaignBySlug :one
SELECT * FROM campaigns
WHERE slug = $1;