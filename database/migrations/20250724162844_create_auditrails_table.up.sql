CREATE TABLE audit_trails (
    id BIGSERIAL PRIMARY KEY,
    ip VARCHAR(45) NOT NULL,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(500) NOT NULL,
    user_id BIGINT,
    payload TEXT,
    response_code INTEGER NOT NULL,
    error_response TEXT,
    response_time_ms INTEGER,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa query
CREATE INDEX idx_audit_trails_created_at ON audit_trails(created_at);
CREATE INDEX idx_audit_trails_user_id ON audit_trails(user_id);
CREATE INDEX idx_audit_trails_response_code ON audit_trails(response_code);
CREATE INDEX idx_audit_trails_ip ON audit_trails(ip); 