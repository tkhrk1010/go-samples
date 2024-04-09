#!/bin/bash

# create tables

ENDPOINT_URL=http://localstack:4566

aws dynamodb create-table --table-name applePurchasing \
    --attribute-definitions AttributeName=id,AttributeType=S AttributeName=PurchacedAt,AttributeType=N \
    --key-schema AttributeName=id,KeyType=HASH AttributeName=PurchacedAt,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url $ENDPOINT_URL

aws dynamodb create-table --table-name orangePurchasing \
    --attribute-definitions AttributeName=id,AttributeType=S AttributeName=PurchacedAt,AttributeType=N \
    --key-schema AttributeName=id,KeyType=HASH AttributeName=PurchacedAt,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url $ENDPOINT_URL

aws dynamodb create-table --table-name bananaPurchasing \
    --attribute-definitions AttributeName=id,AttributeType=S AttributeName=PurchacedAt,AttributeType=N \
    --key-schema AttributeName=id,KeyType=HASH AttributeName=PurchacedAt,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url $ENDPOINT_URL

aws dynamodb create-table --table-name tradeSupportInformation \
    --attribute-definitions AttributeName=id,AttributeType=S AttributeName=executedAt,AttributeType=N AttributeName=value,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH AttributeName=executedAt,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url $ENDPOINT_URL