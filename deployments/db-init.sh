#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER abf WITH password 'qwerty123';
    CREATE DATABASE abf;
    GRANT ALL PRIVILEGES ON DATABASE abf TO abf;
EOSQL
