-- Create enum for user roles
CREATE TYPE user_role AS ENUM ('admin', 'employee');

-- Create users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(150) NOT NULL,
    roles user_role NOT NULL,
    salary DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    created_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT,
    updated_at TIMESTAMP WITH TIME ZONE
);

-- Create index on username for faster lookups
CREATE INDEX idx_users_username ON users(username);

-- Create index on roles for filtering
CREATE INDEX idx_users_roles ON users(roles);

-- Create trigger to automatically update updated_at and updated_by
CREATE OR REPLACE FUNCTION update_updated_columns()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    -- updated_by will be set by application code, not automatically
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_columns 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_columns(); 