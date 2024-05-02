package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	p "github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/persistence"
	"google.golang.org/protobuf/proto"
)

func main() {
	log.Println("start")

	// DynamoDBクライアントの初期化は省略
	client := initializeDynamoDBClient()

	// journalテーブルからデータを読み込む
	journalData, err := scanTable(client, "journal")
	if err != nil {
		log.Fatalf("Failed to scan journal table: %v", err)
	}
	log.Printf("journalData count: %v", len(journalData))

	// snapshotテーブルからデータを読み込む
	snapshotData, err := scanTable(client, "snapshot")
	if err != nil {
		log.Fatalf("Failed to scan snapshot table: %v", err)
	}
	log.Printf("snapshotData count: %v", len(snapshotData))

	// データをデシリアライズして変換し、新しいテーブルに保存
	err = processAndSaveData(client, journalData, "journal_readable")
	if err != nil {
		log.Fatalf("Failed to process and save journal data: %v", err)
	}

	err = processAndSaveData(client, snapshotData, "snapshot_readable")
	if err != nil {
		log.Fatalf("Failed to process and save snapshot data: %v", err)
	}

	log.Println("ETL process completed successfully")
}

func scanTable(client *dynamodb.Client, tableName string) ([]map[string]types.AttributeValue, error) {
	var data []map[string]types.AttributeValue

	paginator := dynamodb.NewScanPaginator(client, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		data = append(data, page.Items...)
	}

	return data, nil
}

func processAndSaveData(client *dynamodb.Client, data []map[string]types.AttributeValue, tableName string) error {
	// データをデシリアライズして変換
	var items []map[string]types.AttributeValue
	for _, item := range data {
		// debug
		fmt.Printf("Item: %+v\n", item)

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
		if err != nil {
			return err
		}
		items = append(items, newItem)

	}

	//debug
	fmt.Printf("Items: %+v\n", items)

	// 変換したデータを新しいテーブルに保存
	var writeReqs []types.WriteRequest
	for _, item := range items {
		writeReqs = append(writeReqs, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		})
	}

	// debug
	fmt.Printf("Write Requests: %+v\n", writeReqs)

	_, err := client.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: writeReqs,
		},
	})
	if err != nil {
		return err
	}

	return nil
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
