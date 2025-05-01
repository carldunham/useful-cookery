-- Create images table
CREATE TABLE IF NOT EXISTS images (
  id TEXT PRIMARY KEY,
  url TEXT NOT NULL,
  alt TEXT,
  width INTEGER,
  height INTEGER
);

-- Create recipe_images junction table
CREATE TABLE IF NOT EXISTS recipe_images (
  recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
  image_id TEXT REFERENCES images(id) ON DELETE CASCADE,
  PRIMARY KEY (recipe_id, image_id)
);

-- Create step_images junction table
CREATE TABLE IF NOT EXISTS step_images (
  step_id TEXT REFERENCES steps(id) ON DELETE CASCADE,
  image_id TEXT REFERENCES images(id) ON DELETE CASCADE,
  PRIMARY KEY (step_id, image_id)
);
