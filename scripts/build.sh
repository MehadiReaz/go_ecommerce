#!/bin/bash

# Build script for the application

set -e

echo "Building E-Commerce API..."

# Build API server
echo "Building API server..."
go build -o bin/api cmd/api/main.go

# Build worker
echo "Building background workers..."
go build -o bin/worker cmd/worker/main.go

# Build seed script
echo "Building seed script..."
go build -o bin/seed scripts/seed_data.go

echo "Build completed successfully!"
echo ""
echo "Binaries created:"
echo "  - bin/api       (API server)"
echo "  - bin/worker    (Background workers)"
echo "  - bin/seed      (Database seeder)"
