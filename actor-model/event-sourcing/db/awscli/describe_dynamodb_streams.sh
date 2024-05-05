#!/bin/bash

#
# Describe DynamoDB Streams
# DynamoDB Streams の設定を確認する。
#
# jqが必要
# https://stedolan.github.io/jq/
# 面倒なので、Macにbrewでinstallしているが、いつかDockerfileに混ぜ込んでもいいかもしれない。
#

cd $(dirname "$0") && pwd

# ストリームのARNを指定
STREAM_ARN="arn:aws:dynamodb:us-east-1:000000000000:table/UserAccountEvent/stream/2024-04-07T04
:04:08.185"

# DynamoDB Streamsの詳細を取得して、最初のシャードIDを抽出
SHARD_ID=$(docker-compose exec -T awscli aws dynamodbstreams describe-stream --stream-arn $STREAM_ARN --endpoint-url http://host.docker.internal:4566 | jq -r '.StreamDescription.Shards[0].ShardId')

echo "Shard ID: $SHARD_ID"

# シャードイテレータを取得
SHARD_ITERATOR=$(docker-compose exec -T awscli aws dynamodbstreams get-shard-iterator --stream-arn $STREAM_ARN --shard-id $SHARD_ID --shard-iterator-type TRIM_HORIZON --endpoint-url http://host.docker.internal:4566 | jq -r '.ShardIterator')

echo "Shard Iterator: $SHARD_ITERATOR"

# シャードイテレータを使用してレコードを取得
docker-compose exec -T awscli aws dynamodbstreams get-records --shard-iterator $SHARD_ITERATOR --endpoint-url http://host.docker.internal:4566
