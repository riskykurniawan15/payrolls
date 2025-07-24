CREATE TABLE reimbursements (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(100) NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    description TEXT,
    created_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT,
    updated_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraint
    CONSTRAINT fk_reimbursements_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_reimbursements_date ON reimbursements(date);
CREATE INDEX idx_reimbursements_user_id ON reimbursements(user_id);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_reimbursements_updated_columns 
    BEFORE UPDATE ON reimbursements 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_columns(); 