#!/bin/bash

# Integration test runner for action tracker

set -e

# Change to project root directory
cd "$(dirname "$0")/.."

echo "Starting integration test environment..."

# Start test database
echo "Starting test database..."
docker-compose -f tests/integration/docker-compose.test.yml up -d test-db

# Wait for database to be healthy
echo "Waiting for database to be healthy..."
sleep 10

# Set environment variables for tests
export DATABASE_URL="postgres://tracker_user:test_password@localhost:5433/action_tracker_test?sslmode=disable"

# Run integration tests
echo "Running integration tests..."
go test -tags=integration -v ./tests/integration/...

TEST_EXIT_CODE=$?

# Clean up
echo "Cleaning up test environment..."
docker-compose -f tests/integration/docker-compose.test.yml down

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "✅ Integration tests passed!"
else
    echo "❌ Integration tests failed!"
    exit $TEST_EXIT_CODE
fi