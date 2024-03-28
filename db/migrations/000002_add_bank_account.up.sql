CREATE TABLE bank_accounts (
  id SERIAL PRIMARY KEY,
  user_id INT,
  account_number VARCHAR(30),
  bank_name VARCHAR(30),
  balance NUMERIC(12),
  currency VARCHAR(10),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_user_id ON bank_accounts (user_id);
