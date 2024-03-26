#!/bin/sh

# setup when you first use
# chmod +x migrate.sh

# 環境変数を読み込む
source .env.dev

# コマンドライン引数をすべて受け取る
COMMAND="$@"

DATABASE_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable"

docker-compose exec migrate migrate -path=/migrations -database "$DATABASE_URL" $COMMAND
