-- Create periods table
CREATE TABLE periods (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status SMALLINT DEFAULT 1,
    user_executable_payroll BIGINT,
    payroll_date TIMESTAMP WITH TIME ZONE,
    created_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT,
    updated_at TIMESTAMP WITH TIME ZONE
);

-- Create index on code for faster lookups
CREATE INDEX idx_periods_code ON periods(code);

-- Create index on start_date for date range queries
CREATE INDEX idx_periods_start ON periods(start_date);

-- Create index on start_date for date range queries
CREATE INDEX idx_periods_end ON periods(end_date);

CREATE TRIGGER update_periods_updated_columns 
    BEFORE UPDATE ON periods 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_columns();
