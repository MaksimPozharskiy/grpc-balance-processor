CREATE TYPE source_t AS ENUM ('game','payment','service');
CREATE TYPE state_t  AS ENUM ('deposit','withdraw');

CREATE TABLE IF NOT EXISTS accounts (
  id         uuid PRIMARY KEY,
  balance    NUMERIC(20,2) NOT NULL DEFAULT 0,
  updated_at timestamptz   NOT NULL DEFAULT now(),
  CONSTRAINT balance_nonneg CHECK (balance >= 0)
);

CREATE TABLE IF NOT EXISTS operations (
  id           bigserial PRIMARY KEY,
  tx_id        text UNIQUE NOT NULL,
  account_id   uuid NOT NULL REFERENCES accounts(id),
  source       source_t NOT NULL,
  state        state_t  NOT NULL,
  amount       NUMERIC(20,2) NOT NULL CHECK (amount > 0),
  created_at   timestamptz NOT NULL DEFAULT now(),
  applied      boolean NOT NULL DEFAULT false,
  canceled_at  timestamptz,
  cancel_note  text
);

CREATE INDEX IF NOT EXISTS idx_ops_account_created ON operations(account_id, created_at);
CREATE INDEX IF NOT EXISTS idx_ops_not_canceled   ON operations(account_id) WHERE canceled_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_ops_applied_open   ON operations(account_id, created_at DESC, id DESC)
  WHERE applied = true AND canceled_at IS NULL;
