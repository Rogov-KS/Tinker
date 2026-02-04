#!/bin/sh
set -e

echo "Starting application..."

# Миграции Prisma отключены (используются Go миграции через db-migrate)
# echo "Applying database migrations..."
# npx prisma migrate deploy

# Проверяем, нужно ли запускать seed
if [ "$SEED_DB" = "true" ] || [ "$SEED_DB" = "1" ]; then
  echo "Running database seed..."
  npx prisma db seed
else
  echo "Skipping seed (set SEED_DB=true to enable)"
fi

# Запускаем приложение
echo "Starting NestJS application..."
exec node dist/src/main.js

