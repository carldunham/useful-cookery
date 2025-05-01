-- Create recipes table
CREATE TABLE IF NOT EXISTS recipes (
  id TEXT PRIMARY KEY,
  original_id TEXT,
  title TEXT NOT NULL,
  description TEXT,
  notes TEXT,
  author_id TEXT REFERENCES users(id),
  cuisine TEXT,
  prep_time INTEGER,
  cook_time INTEGER,
  servings INTEGER,
  difficulty TEXT,
  nutrition_info JSONB,
  tags TEXT [],
  likes INTEGER DEFAULT 0,
  average_rating FLOAT DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
