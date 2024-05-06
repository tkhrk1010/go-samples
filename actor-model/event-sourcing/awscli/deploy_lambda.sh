#!/bin/bash

# 変数の定義
FUNCTION_NAME="read-model-updater"
# handler名はなんでもいい
# https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/golang-handler.html#golang-handler-naming
HANDLER="bootstrap"
RUNTIME="go1.x"
ROLE_ARN="arn:aws:iam::000000000000:role/lambda-role"
ENDPOINT_URL="http://host.docker.internal:4566"

# Lambdaディレクトリに移動
cd ../rmu

# Go依存関係のダウンロード
go mod download

# Lambdaバイナリのビルド
# バイナリ名はbootstrapにする必要がある
# https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/golang-package.html
GOOS=linux GOARCH=amd64 go build -o bootstrap

# ZIPファイルの作成
zip -r ../awscli/function.zip bootstrap

# Lambdaディレクトリから移動
cd ../awscli

# Lambda関数がすでに存在する場合は削除する
docker-compose exec awscli aws lambda delete-function \
  --function-name ${FUNCTION_NAME} \
  --endpoint-url=${ENDPOINT_URL}

# Lambda関数の作成または更新
docker-compose exec awscli aws lambda create-function \
  --function-name ${FUNCTION_NAME} \
  --runtime ${RUNTIME} \
  --handler ${HANDLER} \
  --role ${ROLE_ARN} \
  --zip-file fileb:///awscli/function.zip \
  --endpoint-url=${ENDPOINT_URL} \
  --environment '{ "Variables": {
    "AWS_REGION":"us-east-1",
    "API_HOST":"0.0.0.0",
    "API_PORT":"8080",
    "AWS_REGION_NAME":"us-east-1",
    "AWS_DYNAMODB_ENDPOINT_URL":"http://host.docker.internal:4566",
    "AWS_DYNAMODB_ACCESS_KEY_ID":"x",
    "AWS_DYNAMODB_SECRET_ACCESS_KEY":"x",
    "STREAM_JOURNAL_TABLE_NAME":"journal",
    "STREAM_MAX_ITEM_COUNT":"32",
    "DATABASE_URL":"ceer:ceer@tcp(mysql-local:3306)/ceer"
  } }'

# デプロイ完了のメッセージ
echo "Lambda function ${FUNCTION_NAME} deployed successfully!"
