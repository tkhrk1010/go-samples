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

# ストリームのリストを取得
STREAMS=$(docker-compose exec -T awscli aws dynamodbstreams list-streams --endpoint-url http://host.docker.internal:4566)

# 最初のストリームのARNを取得
STREAM_ARN=$(echo $STREAMS | jq -r '.Streams[0].StreamArn')

echo "Stream ARN: $STREAM_ARN"

# DynamoDB Streamsの詳細を取得して、最初のシャードIDを抽出
SHARD_ID=$(docker-compose exec -T awscli aws dynamodbstreams describe-stream --stream-arn $STREAM_ARN --endpoint-url http://host.docker.internal:4566 | jq -r '.StreamDescription.Shards[0].ShardId')

echo "Shard ID: $SHARD_ID"

# シャードイテレータを取得
SHARD_ITERATOR=$(docker-compose exec -T awscli aws dynamodbstreams get-shard-iterator --stream-arn $STREAM_ARN --shard-id $SHARD_ID --shard-iterator-type TRIM_HORIZON --endpoint-url http://host.docker.internal:4566 | jq -r '.ShardIterator')

echo "Shard Iterator: $SHARD_ITERATOR"

# シャードイテレータを使用してレコードを取得
docker-compose exec -T awscli aws dynamodbstreams get-records --shard-iterator $SHARD_ITERATOR --endpoint-url http://host.docker.internal:4566
