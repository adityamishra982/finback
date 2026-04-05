#!/bin/bash
set -e

echo "Tidying go modules..."
go mod tidy

echo "Building Go API server..."
go build -v -o api.exe ./cmd/api

echo "Build complete! You can run it with ./api.exe"
