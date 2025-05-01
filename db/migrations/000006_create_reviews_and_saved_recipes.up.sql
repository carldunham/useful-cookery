-- Create reviews table
CREATE TABLE IF NOT EXISTS reviews (
  id TEXT PRIMARY KEY,
  recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
  author_id TEXT REFERENCES users(id) ON DELETE CASCADE,
  rating INTEGER NOT NULL,
  comment TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create saved_recipes junction table
CREATE TABLE IF NOT EXISTS saved_recipes (
  user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
  recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, recipe_id)
);
