#!/bin/bash

ENDPOINT_URL="http://host.docker.internal:4566"

cd $(dirname "$0") && pwd

# Create table
# localhostじゃないことに注意。
docker-compose exec awscli aws dynamodb create-table \
    --endpoint-url=$ENDPOINT_URL \
    --table-name journal \
    --attribute-definitions \
        AttributeName=actorName,AttributeType=S \
        AttributeName=eventIndex,AttributeType=N \
    --key-schema \
        AttributeName=actorName,KeyType=HASH \
        AttributeName=eventIndex,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

docker-compose exec awscli aws dynamodb create-table \
    --endpoint-url=$ENDPOINT_URL \
    --table-name snapshot \
    --attribute-definitions \
        AttributeName=actorName,AttributeType=S \
        AttributeName=eventIndex,AttributeType=N \
    --key-schema \
        AttributeName=actorName,KeyType=HASH \
        AttributeName=eventIndex,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
