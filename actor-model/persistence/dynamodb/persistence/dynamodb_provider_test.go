package persistence_test

import (
	"testing"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	p "github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/persistence"
)

func TestNewProviderState(t *testing.T) {
	clinet := InitializeDynamoDBClient()
	ps := p.NewProviderState(clinet)
	if ps == nil {
		t.Error("NewProviderState returned nil")
	}
}

func TestGetState(t *testing.T) {
	clinet := InitializeDynamoDBClient()
	ps := p.NewProviderState(clinet)
	state := ps.GetState()
	if state != ps {
		t.Error("GetState should return the same instance")
	}
}

func TestRestart(t *testing.T) {
	clinet := InitializeDynamoDBClient()
	ps := p.NewProviderState(clinet)
	ps.Restart()
}

func TestGetSnapshotInterval(t *testing.T) {
	clinet := InitializeDynamoDBClient()
	ps := p.NewProviderState(clinet)
	interval := ps.GetSnapshotInterval()
	if interval != 3 {
		t.Errorf("GetSnapshotInterval should return 3, got: %d", interval)
	}
}

func TestGetSnapshot(t *testing.T) {
	tableName := "journal"

	client := InitializeDynamoDBClient()
	ps := p.NewProviderState(client)

	// シードデータの準備
	seedData := []map[string]interface{}{
		{
			"actorName":  "testActor",
			"eventIndex": 1,
			"payload":    encodeEvent(&p.Event{Data: "event1"}),
		},
		{
			"actorName":  "testActor",
			"eventIndex": 2,
			"payload":    encodeEvent(&p.Event{Data: "event2"}),
		},
		{
			"actorName":  "testActor",
			"eventIndex": 3,
			"payload":    encodeEvent(&p.Event{Data: "event3"}),
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
	ps.GetEvents(actorName, eventIndexStart, eventIndexEnd, callback)

	expectedEvents := []*p.Event{
		{Data: "event1"},
		{Data: "event2"},
	}
	assert.Equal(t, len(expectedEvents), len(actualEvents))
	for i, expected := range expectedEvents {
		actual := actualEvents[i].(*p.Event)
		assert.True(t, proto.Equal(expected, actual))
	}

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

func TestPersistSnapshot(t *testing.T) {
	clinet := InitializeDynamoDBClient()
	ps := p.NewProviderState(clinet)
	ps.PersistSnapshot("testActor", 1, &p.Snapshot{})
}

// 呼び出せることだけ確認
func TestDeleteSnapshots(t *testing.T) {
	clinet := InitializeDynamoDBClient()
	ps := p.NewProviderState(clinet)
	ps.DeleteSnapshots("testActor", 1)
}

func TestGetEvents(t *testing.T) {
	clinet := InitializeDynamoDBClient()
	ps := p.NewProviderState(clinet)
	ps.GetEvents("testActor", 1, 2, func(e interface{}) {})
}

func TestPersistEvent(t *testing.T) {
	clinet := InitializeDynamoDBClient()
	ps := p.NewProviderState(clinet)
	ps.PersistEvent("testActor", 1, &p.Event{})
}

// 呼び出せることだけ確認
func TestDeleteEvents(t *testing.T) {
	clinet := InitializeDynamoDBClient()
	ps := p.NewProviderState(clinet)
	ps.DeleteEvents("testActor", 1)
}
