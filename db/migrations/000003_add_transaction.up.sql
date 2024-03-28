CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  transfer_proof_image TEXT,
  sender_id INT,
  recipient_id INT,
  amount NUMERIC(12),
  type VARCHAR(7),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (sender_id) REFERENCES users(id),
  FOREIGN KEY (recipient_id) REFERENCES users(id)
);

CREATE INDEX idx_sender_id ON transactions (sender_id);
CREATE INDEX idx_recipient_id ON transactions (recipient_id);
