CREATE TABLE overtimes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    overtimes_date TIMESTAMP WITH TIME ZONE NOT NULL,
    total_hours_time DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    created_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT,
    updated_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraint
    CONSTRAINT fk_overtimes_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_overtimes_date ON overtimes(overtimes_date);
CREATE INDEX idx_overtimes_user_id ON overtimes(user_id);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_overtimes_updated_columns 
    BEFORE UPDATE ON overtimes 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_columns();
