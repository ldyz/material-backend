-- migrations/003_create_user_projects.sql

CREATE TABLE IF NOT EXISTS user_projects (
  user_id BIGINT NOT NULL,
  project_id BIGINT NOT NULL,
  PRIMARY KEY (user_id, project_id),
  CONSTRAINT fk_up_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_up_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
