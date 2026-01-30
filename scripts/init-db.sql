-- Initialize the database for action tracker
-- This script runs automatically when the TimescaleDB container starts

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Create actions table
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

-- Convert actions table to a hypertable
SELECT create_hypertable('actions', 'timestamp');

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_actions_name_timestamp ON actions (name, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_actions_category_timestamp ON actions (category, timestamp DESC);

-- Create habits table for tracking recurring actions
CREATE TABLE IF NOT EXISTS habits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    category VARCHAR(100),
    target_quantity INTEGER DEFAULT 1,
    unit VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for habits
CREATE INDEX IF NOT EXISTS idx_habits_active ON habits (is_active);
CREATE INDEX IF NOT EXISTS idx_habits_category ON habits (category);