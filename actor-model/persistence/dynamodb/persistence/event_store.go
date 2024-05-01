package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type EventStore struct {
	client *dynamodb.Client
	table  string
}

func NewEventStore(client *dynamodb.Client, table string) *EventStore {
	return &EventStore{
		client: client,
		table:  table,
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
	// TODO: protoreflect.ProtoMessageを実装した構造体に変換する
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
	payload, err := proto.Marshal(event)
	if err != nil {
		panic(err)
	}

	item := map[string]types.AttributeValue{
		"actorName":  &types.AttributeValueMemberS{Value: actorName},
		"eventIndex": &types.AttributeValueMemberN{Value: strconv.Itoa(eventIndex)},
		"payload":    &types.AttributeValueMemberB{Value: payload},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(e.table),
		Item:      item,
	}

	_, err = e.client.PutItem(context.TODO(), input)
	if err != nil {
		panic(err)
	}

}

func (e *EventStore) DeleteEvents(actorName string, inclusiveToIndex int) {}
