package persistence_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	p "github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/persistence"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/proto"

	"github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/persistence/testdata"
)

type mockEvent struct {
	// TODO: Add necessary fields for testing
}

func (m *mockEvent) ProtoReflect() protoreflect.Message {
	// TODO: Implement ProtoReflect method for testing
	return nil
}

func InitializeDynamoDBClient() *dynamodb.Client {
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

func TestEventStore_GetEvents(t *testing.T) {
	tableName := "testEventTable"

	client := InitializeDynamoDBClient()
	eventStore := p.NewEventStore(client, tableName)

	// シードデータの準備
	seedData := []map[string]interface{}{
		{
			"actorName":  "testActor",
			"eventIndex": 1,
			"payload":    []byte(`{"data":"event1"}`),
		},
		{
			"actorName":  "testActor",
			"eventIndex": 2,
			"payload":    []byte(`{"data":"event2"}`),
		},
		{
			"actorName":  "testActor",
			"eventIndex": 3,
			"payload":    []byte(`{"data":"event3"}`),
		},
	}
	for _, item := range seedData {
		av, err := attributevalue.MarshalMap(item)
		assert.NoError(t, err)
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}
		_, err = client.PutItem(context.Background(), input)
		assert.NoError(t, err)
	}

	actorName := "testActor"
	eventIndexStart := 1
	eventIndexEnd := 2

	var actualEvents []interface{}
	callback := func(e interface{}) {
		actualEvents = append(actualEvents, e)
	}

	eventStore.GetEvents(actorName, eventIndexStart, eventIndexEnd, callback)

	expectedEvents := []interface{}{
		map[string]interface{}{"data": "event1"},
		map[string]interface{}{"data": "event2"},
	}
	assert.Equal(t, expectedEvents, actualEvents)

	// クリーンアップ
	for _, item := range seedData {
		key, err := attributevalue.MarshalMap(map[string]interface{}{
			"actorName":  item["actorName"],
			"eventIndex": item["eventIndex"],
		})
		assert.NoError(t, err)
		input := &dynamodb.DeleteItemInput{
			Key:       key,
			TableName: aws.String(tableName),
		}
		_, err = client.DeleteItem(context.Background(), input)
		assert.NoError(t, err)
	}
}

func TestEventStore_PersistEvent(t *testing.T) {
	tableName := "testEventTable"

	client := InitializeDynamoDBClient()
	eventStore := p.NewEventStore(client, tableName)

	actorName := "testActor"
	eventIndex := 1
	eventData := &testdata.TestEvent{Data: "testEvent"}

	eventStore.PersistEvent(actorName, eventIndex, eventData)

	// 保存したイベントの検証
	key, err := attributevalue.MarshalMap(map[string]interface{}{
		"actorName":  actorName,
		"eventIndex": eventIndex,
	})
	assert.NoError(t, err)
	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(tableName),
	}
	result, err := client.GetItem(context.Background(), input)
	assert.NoError(t, err)

	var persistedEvent map[string]interface{}
	err = attributevalue.UnmarshalMap(result.Item, &persistedEvent)
	assert.NoError(t, err)

	assert.Equal(t, actorName, persistedEvent["actorName"])
	// Goでは、UnmarshalMapすると数値はfloat64になるため、intに変換する
	assert.Equal(t, eventIndex, int(persistedEvent["eventIndex"].(float64)))

	var persistedEventData testdata.TestEvent
	err = proto.Unmarshal(persistedEvent["payload"].([]byte), &persistedEventData)
	assert.NoError(t, err)
	assert.Equal(t, eventData.Data, persistedEventData.Data)

	// クリーンアップ
	deleteInput := &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(tableName),
	}
	_, err = client.DeleteItem(context.Background(), deleteInput)
	assert.NoError(t, err)
}

// func TestEventStore_DeleteEvents(t *testing.T) {
// 	eventStore := p.NewEventStore()

// 	// Test case 1: Delete events
// 	eventStore.DeleteEvents("actor1", 10)

// 	// TODO: Verify that the events are deleted correctly from the database

// 	// TODO: Add more test cases for different scenarios
// }
