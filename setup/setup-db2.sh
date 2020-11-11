#!/usr/bin/env bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER drone WITH LOGIN PASSWORD 'drone';
    CREATE DATABASE drone;
    GRANT ALL PRIVILEGES ON DATABASE drone TO drone;
    CREATE USER puppy WITH LOGIN PASSWORD 'puppy';
    CREATE DATABASE puppy;
    GRANT ALL PRIVILEGES ON DATABASE puppy TO puppy;
EOSQL
