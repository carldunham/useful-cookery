-- Remove skill_level column from recipes table
ALTER TABLE
  recipes DROP COLUMN skill_level;

-- Drop the enum type
DROP TYPE skill_level_enum;
