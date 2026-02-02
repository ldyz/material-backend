-- Create stocks table with proper foreign key relationships
CREATE TABLE IF NOT EXISTS stocks (
  id BIGSERIAL PRIMARY KEY,
  material_id BIGINT NOT NULL,
  quantity FLOAT NOT NULL DEFAULT 0,
  safety_stock FLOAT NOT NULL DEFAULT 0,
  location VARCHAR(100),
  unit VARCHAR(20),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE,
  UNIQUE(material_id)
);

CREATE INDEX IF NOT EXISTS idx_stocks_material ON stocks(material_id);
