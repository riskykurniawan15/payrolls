-- Drop trigger
DROP TRIGGER IF EXISTS update_reimbursements_updated_columns ON reimbursements;

-- Drop indexes
DROP INDEX IF EXISTS idx_reimbursements_date;
DROP INDEX IF EXISTS idx_reimbursements_user_id;

-- Drop table
DROP TABLE IF EXISTS reimbursements; 