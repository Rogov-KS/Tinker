#!/bin/sh
set -e

echo "Starting Go backend..."

# Проверяем, нужно ли запускать seed
if [ "$SEED_DB" = "true" ] || [ "$SEED_DB" = "1" ]; then
  echo "Running database seed..."
  ./seed
else
  echo "Skipping seed (set SEED_DB=true to enable)"
fi

# Запускаем приложение
echo "Starting Go application..."
exec ./backend

