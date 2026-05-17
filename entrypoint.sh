#!/bin/sh
set -e

# Wait for database to be ready
until nc -z db 5432; do
  echo "Waiting for database to be ready..."
  sleep 1
done

echo "Database is ready!"

# Run migrations
echo "Running database migrations..."
migrate -path /app/migrations -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@db:5432/${DATABASE_NAME}?sslmode=disable" up

echo "Migrations completed successfully!"

# Run seed
echo "Running database seed..."
./seed

echo "Seed completed successfully!"

# Start the application
exec ./server
