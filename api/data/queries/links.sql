-- name: CreateLink :one
INSERT INTO links (code, original_url, normalized_url, user_id, creator_ip)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetByCode :one
SELECT * FROM links
WHERE code = $1;

-- name: GetByNormalizedURL :one
SELECT * FROM links
WHERE normalized_url = $1;

-- name: CountLinksByCreatorIPSince :one
SELECT count(*) FROM links
WHERE creator_ip = $1 AND created_at >= $2;

-- name: ListLinksByUserID :many
SELECT * FROM links
WHERE user_id = $1
ORDER BY created_at DESC;
