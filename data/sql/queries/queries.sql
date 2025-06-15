-- name: GetProgram :one
SELECT * FROM programs WHERE id = $1;

-- name: GetProgramByExternalID :one
SELECT p.* FROM programs p
JOIN external_sources es ON p.id = es.program_id
WHERE es.source_name = $1 AND es.external_id = $2;

-- name: ListPrograms :many
SELECT * FROM programs
ORDER BY published_at DESC;

-- name: SearchPrograms :many
SELECT * FROM programs
WHERE title ILIKE '%' || $1 || '%'
   OR description ILIKE '%' || $1 || '%'
   OR category ILIKE '%' || $1 || '%'
ORDER BY published_at DESC;

-- name: CreateProgram :one
INSERT INTO programs (title, description, category, language, duration, published_at, source)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateProgram :one
UPDATE programs
SET title = $2, description = $3, category = $4, language = $5, duration = $6, published_at = $7
WHERE id = $1
RETURNING *;

-- name: DeleteProgram :exec
DELETE FROM programs WHERE id = $1;

-- name: CreateExternalSource :one
INSERT INTO external_sources (program_id, source_name, external_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetExternalSources :many
SELECT * FROM external_sources WHERE program_id = $1;
