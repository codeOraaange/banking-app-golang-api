CREATE TABLE balance_income (
  id SERIAL PRIMARY KEY,
  bank_id INT,
  transfer_proof_img TEXT,
  deposited_amount NUMERIC(12),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (bank_id) REFERENCES bank_accounts(id)
);

CREATE INDEX idx_bank_id ON balance_income (bank_id);