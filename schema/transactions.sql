CREATE TABLE IF NOT EXISTS `horseq.transactions` (
  transactions_id STRING, -- UUID
  project_id STRING NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  value_usd FLOAT64 NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);