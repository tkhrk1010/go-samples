#!/bin/bash

ENDPOINT_URL="http://host.docker.internal:4566"
FUNCTION_NAME="read-model-updater"

cd $(dirname "$0") && pwd

# docker-compose exec -T awscli aws logs describe-log-groups \
# 	--endpoint-url=$ENDPOINT_URL

# docker-compose exec -T awscli aws logs describe-log-streams \
# 	--log-group-name "/aws/lambda/$FUNCTION_NAME" \
# 	--endpoint-url=$ENDPOINT_URL


# 最新の20件のログを取得
# 最新の20件のログイベントのメッセージ部分を取得し、タイムスタンプで昇順にソート
# grep '^time='コマンドを使用して、time=から始まるログのみを抽出
# jq -R 'fromjson? | select(type == "object")'コマンドを使用して、抽出されたログメッセージをJSONオブジェクトとして解析し、オブジェクトタイプのログのみを選択
# -Rオプションは、jqに各行をrawな文字列として扱うように指示
# fromjson?は、JSONオブジェクトとして解析できる行のみを選択します。解析できない行はスキップされます。
# select(type == "object")は、オブジェクトタイプのログのみを選択
# docker-compose exec -T awscli aws logs filter-log-events \
#     --endpoint-url=$ENDPOINT_URL \
#     --log-group-name /aws/lambda/$FUNCTION_NAME \
#     --filter-pattern '' \
#     --query 'sort_by(events, &timestamp)[-50:].[message]' \
#     --output text | grep '^time=' | jq -R 'fromjson? | select(type == "object")'

docker-compose exec -T awscli aws logs filter-log-events \
    --endpoint-url=$ENDPOINT_URL \
    --log-group-name /aws/lambda/$FUNCTION_NAME \
    --filter-pattern '' \
    --query 'sort_by(events, &timestamp)[-20:].[message]' \
    --output text | sed -n '/^time=/p' | awk '{gsub(/\\"/," \"")}1' > log.txt
