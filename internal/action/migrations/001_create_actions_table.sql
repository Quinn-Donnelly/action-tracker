-- Action feature migrations
-- 001_create_actions_table.sql

CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    quantity INTEGER DEFAULT 1,
    unit VARCHAR(50),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Convert actions table to a hypertable for time-series optimization
SELECT create_hypertable('actions', 'timestamp');

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_actions_name_timestamp ON actions (name, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_actions_category_timestamp ON actions (category, timestamp DESC);

-- Add basic usage stats view for future analytics
CREATE OR REPLACE VIEW action_stats AS
SELECT 
    name,
    category,
    COUNT(*) as total_actions,
    SUM(quantity) as total_quantity,
    MIN(timestamp) as first_action,
    MAX(timestamp) as last_action
FROM actions 
GROUP BY name, category;