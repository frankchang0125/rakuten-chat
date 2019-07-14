#!/bin/bash

DIR="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_ROOT_USER=root
MYSQL_ROOT_PASSWORD=password

echo "Migrating MySQL..."

until curl -sS "${MYSQL_HOST}:${MYSQL_PORT}" &> /dev/null; do echo "Waiting MySQL..." && sleep 3; done

docker run -v ${DIR}:/migrations --network host migrate/migrate -path=/migrations/sys -database mysql://${MYSQL_ROOT_USER}:${MYSQL_ROOT_PASSWORD}@tcp\(${MYSQL_HOST}:${MYSQL_PORT}\)/sys up

docker run -v ${DIR}:/migrations --network host migrate/migrate -path=/migrations/rakuten -database mysql://${MYSQL_ROOT_USER}:${MYSQL_ROOT_PASSWORD}@tcp\(${MYSQL_HOST}:${MYSQL_PORT}\)/rakuten up
