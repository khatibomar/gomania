-- migrate:up
-- Index for programs.created_at (used in ORDER BY clauses)
CREATE INDEX idx_programs_created_at ON programs (created_at DESC);

-- Index for programs.category_id (foreign key for JOINs)
CREATE INDEX idx_programs_category_id ON programs (category_id);

-- Index for categories.name (used in ORDER BY and search)
CREATE INDEX idx_categories_name ON categories (name);

-- Composite index for programs filtered by category and ordered by created_at
CREATE INDEX idx_programs_category_created_at ON programs (category_id, created_at DESC);

-- Enable pg_trgm extension for text search (if not already enabled)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- GIN indexes for text search on title and description (ILIKE patterns)
CREATE INDEX idx_programs_title_gin ON programs USING gin (title gin_trgm_ops);

CREATE INDEX idx_programs_description_gin ON programs USING gin (description gin_trgm_ops);

-- migrate:down
-- Drop indexes
DROP INDEX IF EXISTS idx_programs_description_gin;

DROP INDEX IF EXISTS idx_programs_title_gin;

DROP INDEX IF EXISTS idx_programs_category_created_at;

DROP INDEX IF EXISTS idx_categories_name;

DROP INDEX IF EXISTS idx_programs_category_id;

DROP INDEX IF EXISTS idx_programs_created_at;
