#!/bin/sh

set -e

echo "run db migration"
echo "=====> $DB_SOURCE"
/app/migrate -path /app/migration -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" -verbose up

echo "start up"
exec "$@"