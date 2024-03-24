.PHONY: up down list-tables describe-table put-item get-item delete-item create-table delete-all-tables

up:
	docker-compose up -d

down:
	docker-compose down

create-table:
	@echo "Creating DynamoDB tables..."
	@docker-compose up -d
	@docker-compose exec localstack aws dynamodb create-table --table-name applePurchasing \
	    --attribute-definitions AttributeName=id,AttributeType=S AttributeName=PurchasedAt,AttributeType=N \
	    --key-schema AttributeName=id,KeyType=HASH AttributeName=PurchasedAt,KeyType=RANGE \
	    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
	    --endpoint-url http://localhost:4566
	@docker-compose exec localstack aws dynamodb create-table --table-name orangePurchasing \
	    --attribute-definitions AttributeName=id,AttributeType=S AttributeName=PurchasedAt,AttributeType=N \
	    --key-schema AttributeName=id,KeyType=HASH AttributeName=PurchasedAt,KeyType=RANGE \
	    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
	    --endpoint-url http://localhost:4566
	@docker-compose exec localstack aws dynamodb create-table --table-name bananaPurchasing \
	    --attribute-definitions AttributeName=id,AttributeType=S AttributeName=PurchasedAt,AttributeType=N \
	    --key-schema AttributeName=id,KeyType=HASH AttributeName=PurchasedAt,KeyType=RANGE \
	    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
	    --endpoint-url http://localhost:4566
	@docker-compose exec localstack aws dynamodb create-table --table-name tradeSupportInformation \
			--attribute-definitions AttributeName=id,AttributeType=S AttributeName=executedAt,AttributeType=N \
			--key-schema AttributeName=id,KeyType=HASH AttributeName=executedAt,KeyType=RANGE \
			--provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
			--endpoint-url http://localhost:4566

	@echo "Tables created successfully."

list-tables:
	docker-compose exec localstack aws dynamodb list-tables --endpoint-url http://localstack:4566

describe-table:
	docker-compose exec localstack aws dynamodb describe-table --table-name $(table) --endpoint-url http://localstack:4566

put-item:
	docker-compose exec localstack aws dynamodb put-item --table-name $(table) --item '$(item)' --endpoint-url http://localstack:4566

get-item:
	docker-compose exec localstack aws dynamodb get-item --table-name $(table) --key '$(key)' --endpoint-url http://localstack:4566

delete-item:
	docker-compose exec localstack aws dynamodb delete-item --table-name $(table) --key '$(key)' --endpoint-url http://localstack:4566


# WARNING: This command will delete all tables in DynamoDB.
delete-all-tables:
	@docker-compose exec localstack sh -c "aws dynamodb list-tables --endpoint-url http://localhost:4566 --output json | jq -r '.TableNames[]' | xargs -I {} aws dynamodb delete-table --table-name {} --endpoint-url http://localhost:4566"
