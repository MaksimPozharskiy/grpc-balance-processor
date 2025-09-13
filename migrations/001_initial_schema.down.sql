DROP INDEX IF EXISTS idx_ops_applied_open;
DROP INDEX IF EXISTS idx_ops_not_canceled;
DROP INDEX IF EXISTS idx_ops_account_created;

DROP TABLE IF EXISTS operations;
DROP TABLE IF EXISTS accounts;

DROP TYPE IF EXISTS state_t;
DROP TYPE IF EXISTS source_t;
