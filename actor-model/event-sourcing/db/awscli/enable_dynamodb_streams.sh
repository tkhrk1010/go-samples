#!/bin/bash

# Subscribe DynamoDB Streams

cd $(dirname "$0") && pwd

# KEYS_ONLY：アイテムのキーのみをストリームに含めます。
# NEW_IMAGE：変更後のアイテム全体をストリームに含めます。
# OLD_IMAGE：変更前のアイテム全体をストリームに含めます。
# NEW_AND_OLD_IMAGES：変更前後のアイテム全体をストリームに含めます。
docker-compose exec awscli aws dynamodb update-table \
  --endpoint-url http://host.docker.internal:4566 \
  --table-name journal \
  --stream-specification StreamEnabled=true,StreamViewType=NEW_IMAGE

docker-compose exec awscli aws dynamodb update-table \
  --endpoint-url http://host.docker.internal:4566 \
  --table-name snapshot \
  --stream-specification StreamEnabled=true,StreamViewType=NEW_IMAGE
