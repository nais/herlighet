#!/usr/bin/env bash

set -e

go build

cleanup() {
    kill $GLORYPID
}

trap cleanup EXIT

# So that we don't get the stupid question 
# about whether we want to allow incoming
# connections
if [[ $(uname -s) == "Darwin" ]]; then 
    codesign --force --deep --sign - ./herlighet 
fi
./herlighet & 
GLORYPID=$!

echo "Waiting for herlighet..."
while ! nc -z localhost 5432 > /dev/null; do
  sleep 0.1 # wait for 1/10 of the second before check again
done

echo "herlighet open for business"
(PGPASSWORD=handler PGSSLMODE=disable psql --host 127.0.0.1 -P footer --dbname handler --user handler -c 'SELECT CURRENT_DATABASE() AS "db", PI() AS "Value of PI";')
(PGPASSWORD=puppy PGSSLMODE=disable psql --host 127.0.0.1 -P footer --dbname puppy --user puppy -c "SELECT CURRENT_DATABASE() AS "db", 'Do not google; Republican inside' AS \"Santorum\";")
echo
echo "This should give an error saying connections to this database are explicitly not allowed:"
echo
(PGPASSWORD=pony PGSSLMODE=disable psql --host 127.0.0.1 -P footer --dbname pony --user pony -c "SELECT CURRENT_DATABASE() AS "db", 'The UN-stable!' AS \"Where the three-legged pony lives\";" || exit 0)
echo
echo "This should give an error saying the database is not known:"
echo
(PGPASSWORD=gimp PGSSLMODE=disable psql --host 127.0.0.1 -P footer --dbname gimp --user gimp -c "SELECT CURRENT_DATABASE() AS "db", 'Best wake him up!' AS \"The gimp is sleeping\";" || exit 0)
