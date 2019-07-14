#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

DIR="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

ES_HOST=localhost
ES_PORT=9200

echo "Migrating Elasticsearch..."

until curl -sS "http://${ES_HOST}:${ES_PORT}" &> /dev/null; do echo "Waiting Elasticsearch..." && sleep 3; done

curl -XPUT -H "Content-Type: application/json" --data "@${DIR}/1563031407/messages_template.json" "http://${ES_HOST}:${ES_PORT}/_template/messages_template"