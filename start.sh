#!/bin/sh

set -e

cmd="$@"

# wait for postgres to be ready
until PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"

# make migrations
goose -dir ./migrations up

# run tests
# go test -v ./...

exec $cmd