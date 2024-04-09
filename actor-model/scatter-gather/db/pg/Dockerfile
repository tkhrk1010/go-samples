FROM postgres:latest

# PostgreSQLの公式イメージは、/docker-entrypoint-initdb.d/ ディレクトリに配置されたスクリプトを初期化時に実行する
COPY ./init.sql /docker-entrypoint-initdb.d/
