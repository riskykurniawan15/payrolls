-- Drop trigger
DROP TRIGGER IF EXISTS update_overtimes_updated_columns ON overtimes;

-- Drop indexes
DROP INDEX IF EXISTS idx_overtimes_date;
DROP INDEX IF EXISTS idx_overtimes_user_id;

-- Drop table
DROP TABLE IF EXISTS overtimes;
