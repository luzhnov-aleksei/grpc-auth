BEGIN;
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email VARCHAR(255) UNIQUE NOT NULL,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       first_name VARCHAR(100),
                       last_name VARCHAR(100),
                       last_login_at TIMESTAMP WITH TIME ZONE,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_username ON users(username);
COMMIT;