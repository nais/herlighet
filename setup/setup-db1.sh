#!/usr/bin/env bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER herlighet WITH LOGIN PASSWORD 'herlighet';
    CREATE DATABASE herlighet;
    \c herlighet;
    CREATE TABLE databases (
        dbname VARCHAR PRIMARY KEY,
        hostname VARCHAR NOT NULL,
        naisdevice BOOLEAN
    );
    ALTER TABLE databases OWNER TO herlighet;
    INSERT INTO public.databases(dbname, hostname, naisdevice) VALUES ('handler', 'localhost:5433', true);
    INSERT INTO public.databases(dbname, hostname, naisdevice) VALUES ('daddy', 'localhost:5433', true);
    INSERT INTO public.databases(dbname, hostname, naisdevice) VALUES ('drone', 'localhost:5434', true);
    INSERT INTO public.databases(dbname, hostname, naisdevice) VALUES ('puppy', 'localhost:5434', true);
    INSERT INTO public.databases(dbname, hostname, naisdevice) VALUES ('pony', 'localhost:5434', false);
    GRANT ALL PRIVILEGES ON DATABASE herlighet TO herlighet;

    CREATE USER handler WITH LOGIN PASSWORD 'handler';
    CREATE DATABASE handler;
    GRANT ALL PRIVILEGES ON DATABASE handler TO handler;
    CREATE USER daddy WITH LOGIN PASSWORD 'daddy';
    CREATE DATABASE daddy;
    GRANT ALL PRIVILEGES ON DATABASE daddy TO daddy;
EOSQL
