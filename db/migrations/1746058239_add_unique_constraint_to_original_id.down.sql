-- Migration Down
-- Drop the unique index on the original_id column
DROP INDEX IF EXISTS recipes_original_id_unique_idx;
