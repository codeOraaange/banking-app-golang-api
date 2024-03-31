CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  password VARCHAR(100),
  email VARCHAR(50) UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_password ON users (password);
CREATE INDEX idx_users_email ON users (email);