-- Create enum type for skill_level
CREATE TYPE skill_level_enum AS ENUM ('BEGINNER', 'INTERMEDIATE', 'ADVANCED');

-- Add skill_level column to recipes table
ALTER TABLE
  recipes
ADD
  COLUMN skill_level skill_level_enum;

-- Update existing recipes to set skill_level based on difficulty
-- Map 'Easy' to 'BEGINNER', 'Medium' to 'INTERMEDIATE', 'Hard' to 'ADVANCED'
UPDATE
  recipes
SET
  skill_level = 'BEGINNER' :: skill_level_enum
WHERE
  difficulty = 'Easy';

UPDATE
  recipes
SET
  skill_level = 'INTERMEDIATE' :: skill_level_enum
WHERE
  difficulty = 'Medium';

UPDATE
  recipes
SET
  skill_level = 'ADVANCED' :: skill_level_enum
WHERE
  difficulty = 'Hard';
