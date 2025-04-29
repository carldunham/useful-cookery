-- Create ingredient_units table
CREATE TABLE IF NOT EXISTS ingredient_units (
  id TEXT PRIMARY KEY,
  ingredient_id TEXT REFERENCES ingredients(id) ON DELETE CASCADE,
  system TEXT NOT NULL,
  value FLOAT NOT NULL,
  unit TEXT NOT NULL,
  is_main BOOLEAN DEFAULT FALSE
);

-- Add index for faster lookups
CREATE INDEX idx_ingredient_units_ingredient_id ON ingredient_units(ingredient_id);
