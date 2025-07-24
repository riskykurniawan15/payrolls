-- Drop trigger
DROP TRIGGER IF EXISTS update_attendances_updated_columns ON attendances;

-- Drop indexes
DROP INDEX IF EXISTS idx_attendances_user_id;
DROP INDEX IF EXISTS idx_attendances_check_out_date;
DROP INDEX IF EXISTS idx_attendances_check_in_date;

-- Drop table
DROP TABLE IF EXISTS attendances;
