-- Drop trigger first
DROP TRIGGER IF EXISTS update_users_updated_columns ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_columns();

-- Drop indexes
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_roles;

-- Drop table
DROP TABLE IF EXISTS users;

-- Drop enum type
DROP TYPE IF EXISTS user_role; 