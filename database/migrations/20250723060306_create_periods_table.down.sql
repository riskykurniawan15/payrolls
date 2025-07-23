-- Drop trigger first
DROP TRIGGER IF EXISTS update_periods_updated_columns ON periods;

-- Drop indexes
DROP INDEX IF EXISTS idx_periods_code;
DROP INDEX IF EXISTS idx_periods_start;
DROP INDEX IF EXISTS idx_periods_end;

-- Drop periods table
DROP TABLE IF EXISTS periods;
