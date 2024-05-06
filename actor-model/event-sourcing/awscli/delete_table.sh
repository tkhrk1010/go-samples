#!/bin/bash

ENDPOINT_URL="http://host.docker.internal:4566"

cd $(dirname "$0") && pwd

# Delete table
docker-compose exec awscli aws dynamodb delete-table \
    --endpoint-url=$ENDPOINT_URL \
    --table-name journal
docker-compose exec awscli aws dynamodb delete-table \
    --endpoint-url=$ENDPOINT_URL \
    --table-name snapshot
