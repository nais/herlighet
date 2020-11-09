#!/usr/bin/env bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER handler WITH LOGIN PASSWORD 'handler';
    CREATE DATABASE handler;
    GRANT ALL PRIVILEGES ON DATABASE handler TO handler;
    CREATE USER daddy WITH LOGIN PASSWORD 'daddy';
    CREATE DATABASE daddy;
    GRANT ALL PRIVILEGES ON DATABASE daddy TO daddy;
EOSQL
