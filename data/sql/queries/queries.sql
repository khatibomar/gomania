-- name: GetProgram :one
SELECT
    p.id,
    p.title,
    p.description,
    p.language,
    p.duration,
    p.category_id,
    c.name as category_name
FROM programs p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.id = $1;

-- name: ListPrograms :many
SELECT
    p.id,
    p.title,
    p.description,
    p.language,
    p.duration,
    p.category_id,
    c.name as category_name
FROM programs p
LEFT JOIN categories c ON p.category_id = c.id
ORDER BY p.created_at DESC;

-- name: SearchPrograms :many
SELECT
    p.id,
    p.title,
    p.description,
    p.language,
    p.duration,
    p.category_id,
    c.name as category_name
FROM programs p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.title ILIKE '%' || $1 || '%'
   OR p.description ILIKE '%' || $1 || '%'
   OR c.name ILIKE '%' || $1 || '%'
ORDER BY p.created_at DESC;

-- name: CreateProgram :one
INSERT INTO programs (title, description, category_id, language, duration)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, title, description, language, duration;

-- name: UpdateProgram :one
UPDATE programs
SET
    title = $2,
    description = $3,
    category_id = $4,
    language = $5,
    duration = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, title, description, language, duration;

-- name: DeleteProgram :exec
DELETE FROM programs WHERE id = $1;

-- name: GetCategories :many
SELECT id, name
FROM categories
ORDER BY name;

-- name: CreateCategory :one
INSERT INTO categories (name)
VALUES ($1)
RETURNING id, name;

-- name: GetProgramsByCategory :many
SELECT
    p.id,
    p.title,
    p.description,
    p.language,
    p.duration,
    p.category_id,
    c.name as category_name
FROM programs p
JOIN categories c ON p.category_id = c.id
WHERE c.id = $1
ORDER BY p.created_at DESC;
