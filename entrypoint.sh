#!/bin/sh
set -e

# Wait for database to be ready
echo "Waiting for database to be ready..."
until pg_isready -d "$DATABASE_URL" -q; do
  sleep 2
done

echo "Database is ready!"

# Run migrations
echo "Running database migrations..."
migrate -path /app/migrations -database "$DATABASE_URL" up

echo "Migrations completed successfully!"

# Run seed
echo "Running database seed..."
./seed

echo "Seed completed successfully!"

# Start the application
exec ./server