#!/bin/bash

# 手動でlambdaをkickしてtestしたいときに使う

ENDPOINT_URL="http://host.docker.internal:4566"
FUNCTION_NAME="read-model-updater"

cd $(dirname "$0") && pwd

# lambdaが存在することを確認
# docker-compose exec -T awscli aws lambda get-function \
# 	--endpoint-url=$ENDPOINT_URL \
# 	--function-name $FUNCTION_NAME \
# 	--output text

docker-compose exec awscli aws lambda invoke \
	--cli-binary-format raw-in-base64-out \
	--function-name $FUNCTION_NAME \
	--payload file:///awscli/example-dynamodb-event.json \
	--endpoint-url=$ENDPOINT_URL \
	/awscli/output.txt
