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

-- name: GetTotalUserCampaigns :one
SELECT COUNT(*) AS total
FROM campaigns
WHERE 
	user_id = $1 AND
	deleted_at IS NULL AND
	title ILIKE '%' || sqlc.arg(title)::text || '%' AND
	status = sqlc.arg(status)::integer;

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
SELECT campaigns.id, campaigns.title, campaigns.description, campaigns.slug, campaigns.target_amount, campaigns.current_amount, campaigns.start_date, campaigns.end_date,
	users.name as user_name, users.email as user_email,
	CASE 
		WHEN campaigns.current_amount = 0 THEN 0 
		ELSE campaigns.target_amount / campaigns.current_amount 
	END::DECIMAL(10, 2) AS progress
FROM campaigns
JOIN users ON campaigns.user_id = users.id
WHERE campaigns.slug = $1;

-- name: GetCampaigns :many
SELECT id, title, slug,
		current_amount, target_amount,
	   CASE 
		   WHEN current_amount = 0 THEN 0 
		   ELSE target_amount / current_amount 
	   END::DECIMAL(10, 2) AS progress, 
	   start_date, end_date,
	   CASE
	   	   	WHEN status = 0 THEN 'Draft'
	   	   	WHEN status = 1 THEN 'Active'
	   	   	WHEN status = 2 THEN 'Completed'
	   	   	WHEN status = 3 THEN 'Cancelled'
	   	   ELSE 'Unknown'
	   END AS status
FROM campaigns
WHERE 
	deleted_at IS NULL AND
	status = 1 AND
	start_date <= CURRENT_TIMESTAMP AND
	end_date >= CURRENT_TIMESTAMP
ORDER BY start_date DESC
LIMIT $1 OFFSET $2;

-- name: GetTotalCampaigns :one
SELECT COUNT(*) AS total
FROM campaigns
WHERE 
	deleted_at IS NULL AND
	status = 1 AND
	start_date <= CURRENT_TIMESTAMP AND
	end_date >= CURRENT_TIMESTAMP;

-- name: FindCampaignsBySlugForUpdate :one
SELECT id, user_id FROM campaigns
WHERE slug = $1 AND deleted_at IS NULL
FOR UPDATE;

-- name: Donate :exec
UPDATE campaigns
SET current_amount = current_amount + sqlc.arg(amount)::numeric	
WHERE id = $1 AND deleted_at IS NULL;