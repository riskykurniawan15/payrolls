CREATE TABLE attendances (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    check_in_date TIMESTAMP WITH TIME ZONE NOT NULL,
    check_out_date TIMESTAMP WITH TIME ZONE NULL,
    created_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT,
    updated_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraint
    CONSTRAINT fk_attendances_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_attendances_check_in_date ON attendances(check_in_date);
CREATE INDEX idx_attendances_check_out_date ON attendances(check_out_date);
CREATE INDEX idx_attendances_user_id ON attendances(user_id);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_attendances_updated_columns 
    BEFORE UPDATE ON attendances 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_columns();
