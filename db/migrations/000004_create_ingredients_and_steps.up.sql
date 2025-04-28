-- Create ingredients table
CREATE TABLE IF NOT EXISTS ingredients (
  id TEXT PRIMARY KEY,
  recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  quantity FLOAT,
  unit TEXT,
  preparation TEXT,
  substitutes TEXT [],
  is_optional BOOLEAN DEFAULT FALSE
);

-- Create steps table
CREATE TABLE IF NOT EXISTS steps (
  id TEXT PRIMARY KEY,
  recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
  order_index INTEGER NOT NULL,
  description TEXT NOT NULL,
  time_estimate INTEGER
);
