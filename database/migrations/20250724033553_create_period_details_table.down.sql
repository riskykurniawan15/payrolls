-- Drop trigger
DROP TRIGGER IF EXISTS update_period_details_updated_columns ON period_details;

-- Drop indexes
DROP INDEX IF EXISTS idx_period_details_periods_id;
DROP INDEX IF EXISTS idx_period_details_user_id;
DROP INDEX IF EXISTS idx_period_details_periods_user;

-- Drop table
DROP TABLE IF EXISTS period_details; 