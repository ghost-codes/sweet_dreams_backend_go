#!/bin/sh

set -e
FILE=/app/.env
if [ ! -f "$FILE" ]; then
    touch .env
fi

echo "run db migration"
echo "=====> $DB_SOURCE"
/app/migrate -path /app/migrations -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" -verbose up


echo "start up"
exec "$@"