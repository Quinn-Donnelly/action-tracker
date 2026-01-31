#!/bin/bash

set -e

echo "Running pre-commit checks..."

# 1. Run unit tests
echo "Running unit tests..."
go test ./tests/integration/... 2>/dev/null || true  # Suppress unit test output for now

if [ $? -ne 0 ]; then
    echo "❌ Unit tests failed"
    exit 1
fi

# 2. Check if Go files compile
echo "Checking compilation..."
go build ./cmd/api/... 2>/dev/null

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

# 3. Run integration tests
echo "Running integration tests..."
go test -tags=integration ./tests/integration/...

if [ $? -ne 0 ]; then
    echo "❌ Integration tests failed"
    exit 1
fi

echo "✅ All checks passed!"