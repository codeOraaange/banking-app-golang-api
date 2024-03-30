CREATE TABLE balance_outcome (
  id SERIAL PRIMARY KEY,
  sender_id INT,
  recipient_id INT,
  credited_amount NUMERIC(12),
  from_currency VARCHAR(5),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (sender_id) REFERENCES bank_accounts(id),
  FOREIGN KEY (recipient_id) REFERENCES bank_accounts(id)
);

CREATE INDEX idx_sender_id ON balance_outcome (sender_id);
CREATE INDEX idx_recipient_id ON balance_outcome (recipient_id);