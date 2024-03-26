# migrate.shスクリプトのパスを指定
MIGRATE_SCRIPT=./migrate.sh

.PHONY: up down status create

# マイグレーションを1つ進める
up:
	@${MIGRATE_SCRIPT} up

# マイグレーションを1つ戻す
down:
	@${MIGRATE_SCRIPT} down 1

# データベースの現在のマイグレーション状態を確認
status:
	@${MIGRATE_SCRIPT} version

# 新しいマイグレーションファイルを作成
# 使用例: make create name=create_users_table
create:
	@${MIGRATE_SCRIPT} create -ext sql -dir /migrations -seq $(name)

# マイグレーションを指定したバージョンに強制的に移動
# 使用例: make force version=1
force:
	@${MIGRATE_SCRIPT} force $(version)