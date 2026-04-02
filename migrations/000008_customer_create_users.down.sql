DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS idx_users_phone;

DROP TRIGGER IF EXISTS set_users_updated_at ON users;
