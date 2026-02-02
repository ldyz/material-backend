-- migrations/002_create_projects_tables.sql

CREATE TABLE IF NOT EXISTS projects (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  code VARCHAR(50),
  location VARCHAR(100),
  start_date DATE,
  end_date DATE,
  description VARCHAR(500),
  manager VARCHAR(50),
  contact VARCHAR(50),
  budget DECIMAL(15,2) DEFAULT 0,
  status VARCHAR(20) DEFAULT 'planning',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_status ON projects(status);
CREATE INDEX idx_manager ON projects(manager);
