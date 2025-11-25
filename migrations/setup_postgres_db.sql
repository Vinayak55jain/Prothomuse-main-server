-- Run this script in PostgreSQL to set up the postgres database for prothomuse
-- Connect to postgres database first: psql -U postgres -d postgres

-- Drop existing tables if they exist (CAREFUL: this will delete data!)
-- DROP TABLE IF EXISTS metrics CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_api_key ON users(api_key);

-- Create metrics table
CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    project_id VARCHAR(255) NOT NULL,
    route VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    status_code INT NOT NULL,
    response_time BIGINT NOT NULL,
    timestamp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes on metrics
CREATE INDEX IF NOT EXISTS idx_project_id ON metrics(project_id);
CREATE INDEX IF NOT EXISTS idx_timestamp ON metrics(timestamp);

-- Verify tables were created
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'public' ORDER BY table_name;
