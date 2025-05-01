-- Migration Up
-- Add a unique constraint to the original_id column in the recipes table
-- First, handle any existing NULL values
UPDATE
  recipes
SET
  original_id = ''
WHERE
  original_id IS NULL;

-- Create a temporary table to identify duplicates
CREATE TEMP TABLE duplicate_original_ids AS
SELECT
  original_id,
  array_agg(id) AS recipe_ids
FROM
  recipes
WHERE
  original_id != ''
GROUP BY
  original_id
HAVING
  COUNT(*) > 1;

-- Update duplicate original_ids to make them unique by appending a suffix
UPDATE
  recipes r
SET
  original_id = r.original_id || '-' || idx
FROM
  (
    SELECT
      d.original_id,
      recipe_ids,
      unnest(recipe_ids) AS recipe_id,
      generate_subscripts(recipe_ids, 1) AS idx
    FROM
      duplicate_original_ids d
  ) AS subq
WHERE
  r.id = subq.recipe_id
  AND idx > 1;

-- Create a partial unique index on the original_id column
-- This will exclude empty strings
CREATE UNIQUE INDEX recipes_original_id_unique_idx ON recipes (original_id)
WHERE
  original_id != '';

-- Drop the temporary table
DROP TABLE duplicate_original_ids;
