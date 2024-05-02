package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	p "github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/persistence"
	"google.golang.org/protobuf/proto"
)

func main() {
	log.Println("start")

	//
	// DynamoDBクライアントの初期化
	client := initializeDynamoDBClient()

	//
	// 出力先tableの作成
	err := createTableIfNotExists(client, "journal_readable")
	if err != nil {
		log.Fatalf("Failed to create journal_readable table: %v", err)
	}

	err = createTableIfNotExists(client, "snapshot_readable")
	if err != nil {
		log.Fatalf("Failed to create snapshot_readable table: %v", err)
	}

	//
	// 変換処理
	batchSize := 10

	// journalテーブルからデータを読み込み、処理する
	err = processBatchData(client, "journal", "journal_readable", batchSize)
	if err != nil {
		log.Fatalf("Failed to process journal data: %v", err)
	}

	// snapshotテーブルからデータを読み込み、処理する
	err = processBatchData(client, "snapshot", "snapshot_readable", batchSize)
	if err != nil {
		log.Fatalf("Failed to process snapshot data: %v", err)
	}

	log.Println("ETL process completed successfully")
}

func scanTableWithLimit(client *dynamodb.Client, tableName string, startKey map[string]types.AttributeValue, limit int) ([]map[string]types.AttributeValue, map[string]types.AttributeValue, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int32(int32(limit)),
	}
	if startKey != nil {
		input.ExclusiveStartKey = startKey
	}

	output, err := client.Scan(context.TODO(), input)
	if err != nil {
		return nil, nil, err
	}

	return output.Items, output.LastEvaluatedKey, nil
}

func processBatchData(client *dynamodb.Client, srcTableName, dstTableName string, batchSize int) error {
	var startKey map[string]types.AttributeValue

	for {
		// 一定サイズのデータを読み込む
		data, lastEvaluatedKey, err := scanTableWithLimit(client, srcTableName, startKey, batchSize)
		if err != nil {
			return fmt.Errorf("Failed to scan %s table: %v", srcTableName, err)
		}

		// 読み込んだデータを処理し、保存する
		err = processAndSaveData(client, data, dstTableName, batchSize)
		if err != nil {
			return fmt.Errorf("Failed to process and save %s data: %v", srcTableName, err)
		}

		if lastEvaluatedKey == nil {
			break
		}
		startKey = lastEvaluatedKey
	}

	return nil
}

func processAndSaveData(client *dynamodb.Client, data []map[string]types.AttributeValue, tableName string, batchSize int) error {
	var writeReqs []types.WriteRequest

	for _, item := range data {
		payload, ok := item["payload"].(*types.AttributeValueMemberB)
		if !ok {
			return fmt.Errorf("payload is not a binary type")
		}
		event := &p.Event{}
		err := proto.Unmarshal(payload.Value, event)
		if err != nil {
			return err
		}

		actorName := item["actorName"].(*types.AttributeValueMemberS).Value
		eventIndex := item["eventIndex"].(*types.AttributeValueMemberN).Value

		// 変換したデータを新しいitemに追加
		newItem := map[string]types.AttributeValue{
			"actorName":  &types.AttributeValueMemberS{Value: actorName},
			"eventIndex": &types.AttributeValueMemberN{Value: eventIndex},
			"payload":    &types.AttributeValueMemberS{Value: event.String()},
		}

		writeReqs = append(writeReqs, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: newItem,
			},
		})

		// バッチサイズに達したら書き込む
		if len(writeReqs) == batchSize {
			err := batchWriteItems(client, tableName, writeReqs)
			if err != nil {
				return err
			}
			writeReqs = []types.WriteRequest{}
		}
	}

	// 残りのリクエストを書き込む
	if len(writeReqs) > 0 {
		err := batchWriteItems(client, tableName, writeReqs)
		if err != nil {
			return err
		}
	}

	return nil
}

func batchWriteItems(client *dynamodb.Client, tableName string, writeReqs []types.WriteRequest) error {
	_, err := client.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: writeReqs,
		},
	})
	return err
}

func initializeDynamoDBClient() *dynamodb.Client {
	ctx := context.TODO()

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:4566",
				SigningRegion: "us-east-1",
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	return dynamodb.NewFromConfig(cfg)
}

func createTableIfNotExists(client *dynamodb.Client, tableName string) error {
	_, err := client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		var notFoundErr *types.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			// テーブルが存在しない場合は作成する
			_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
				AttributeDefinitions: []types.AttributeDefinition{
					{
						AttributeName: aws.String("actorName"),
						AttributeType: types.ScalarAttributeTypeS,
					},
					{
						AttributeName: aws.String("eventIndex"),
						AttributeType: types.ScalarAttributeTypeN,
					},
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("actorName"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("eventIndex"),
						KeyType:       types.KeyTypeRange,
					},
				},
				TableName: aws.String(tableName),
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(3),
					WriteCapacityUnits: aws.Int64(3),
				},
			})
			if err != nil {
				return err
			}
			// テーブルが作成されるまで待つ
			err = waitUntilTableExists(client, tableName)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func waitUntilTableExists(client *dynamodb.Client, tableName string) error {
	waiter := dynamodb.NewTableExistsWaiter(client)
	err := waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}, 5*time.Minute)
	if err != nil {
		return err
	}
	return nil
}
