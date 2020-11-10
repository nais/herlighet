#!/usr/bin/env bash

set -e

go build

cleanup() {
    kill $GLORYPID
}

trap cleanup EXIT

./herlighet & 
GLORYPID=$!

echo "Waiting for herlighet to open..."
while ! nc -z localhost 5432 > /dev/null; do
  sleep 0.1 # wait for 1/10 of the second before check again
done

echo "herlighet open for business"
(PGPASSWORD=handler PGSSLMODE=disable psql --host 127.0.0.1 --dbname handler --user handler -c 'SELECT PI() AS "Value of PI";')
(PGPASSWORD=puppy PGSSLMODE=disable psql --host 127.0.0.1 --dbname puppy --user puppy -c "SELECT 'Do not google; Republican inside' AS \"Santorum\";")
(PGPASSWORD=pony PGSSLMODE=disable psql --host 127.0.0.1 --dbname pony --user pony -c "SELECT 'The UN-stable!' AS \"Where the three-legged pony lives\";") || exit 0
