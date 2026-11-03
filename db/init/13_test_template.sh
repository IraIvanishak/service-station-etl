#!/bin/sh
# Snapshot the freshly initialised database as a template. Integration tests
# clone it (CREATE DATABASE ... TEMPLATE) so they never touch the dev data.
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres <<SQL
CREATE DATABASE service_station_template TEMPLATE "$POSTGRES_DB";
ALTER DATABASE service_station_template WITH IS_TEMPLATE true ALLOW_CONNECTIONS false;
SQL
