-- name: CreateLink :one
INSERT INTO links (code, original_url, normalized_url)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetByCode :one
SELECT * FROM links
WHERE code = $1;

-- name: GetByNormalizedURL :one
SELECT * FROM links
WHERE normalized_url = $1;
