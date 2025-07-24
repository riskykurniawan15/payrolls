CREATE TABLE period_details (
    id BIGSERIAL PRIMARY KEY,
    periods_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    daily_rate DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    total_working INTEGER NOT NULL DEFAULT 0,
    amount_salary DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    overtime JSONB,
    amount_overtime DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    reimbursement JSONB,
    amount_reimbursement DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    take_home_pay DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    created_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT,
    updated_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraints
    CONSTRAINT fk_period_details_periods_id FOREIGN KEY (periods_id) REFERENCES periods(id) ON DELETE CASCADE,
    CONSTRAINT fk_period_details_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_period_details_periods_id ON period_details(periods_id);
CREATE INDEX idx_period_details_user_id ON period_details(user_id);
CREATE INDEX idx_period_details_periods_user ON period_details(periods_id, user_id);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_period_details_updated_columns
    BEFORE UPDATE ON period_details
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_columns(); 