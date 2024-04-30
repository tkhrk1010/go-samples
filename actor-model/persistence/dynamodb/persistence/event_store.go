package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	// "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type EventStore struct {
	client   *dynamodb.Client
	table    string
}

func NewEventStore(client *dynamodb.Client, table string) *EventStore {
	return &EventStore{
		client:   client,
		table:    table,
	}
}

func (e *EventStore) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{})) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(e.table),
		KeyConditionExpression: aws.String("actorName = :actorName AND eventIndex BETWEEN :start AND :end"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":actorName": &types.AttributeValueMemberS{Value: actorName},
			":start":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", eventIndexStart)},
			":end":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", eventIndexEnd)},
		},
	}

	resp, err := e.client.Query(context.Background(), input)
	if err != nil {
		panic(err)
	}
	for _, item := range resp.Items {
		var event map[string]interface{}
		err := json.Unmarshal([]byte(item["payload"].(*types.AttributeValueMemberB).Value), &event)
		if err != nil {
			panic(err)
		}
		callback(event)
	
	}
}

func (e *EventStore) PersistEvent(actorName string, eventIndex int, event protoreflect.ProtoMessage) {
	// TODO: Implement persisting the event to the database
	// 1. Convert the event to a format suitable for storing in the database (e.g., JSON, binary)
	// 2. Store the converted event in the database, associating it with the actorName and eventIndex
}

func (e *EventStore) DeleteEvents(actorName string, inclusiveToIndex int) {
	// TODO: Implement deleting events from the database
	// 1. Delete all events for the given actorName up to and including the specified inclusiveToIndex
}