-- migrations/004_create_materials_table.sql

CREATE TABLE IF NOT EXISTS materials (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(50) UNIQUE,
  name VARCHAR(50) NOT NULL,
  specification VARCHAR(100),
  unit VARCHAR(20),
  price DECIMAL(15,2) DEFAULT 0,
  description TEXT,
  category VARCHAR(50),
  quantity INT DEFAULT 0,
  project_id BIGINT,
  material VARCHAR(50),
  spec VARCHAR(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_materials_project ON materials(project_id);
CREATE INDEX IF NOT EXISTS idx_materials_category ON materials(category);
CREATE INDEX IF NOT EXISTS idx_materials_name ON materials(name);
