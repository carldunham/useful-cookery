-- Create recipe_categories junction table
CREATE TABLE IF NOT EXISTS recipe_categories (
  recipe_id TEXT REFERENCES recipes(id) ON DELETE CASCADE,
  category_id TEXT REFERENCES categories(id) ON DELETE CASCADE,
  PRIMARY KEY (recipe_id, category_id)
);
