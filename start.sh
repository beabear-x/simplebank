#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migration -database "mysql://root:root@tcp(mysql:3306)/simple_bank" -verbose up

echo "start the app"
exec "$@"