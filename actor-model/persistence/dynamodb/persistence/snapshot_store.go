package persistence

import (
	"fmt"
	"strconv"

	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type snapshotEntry struct {
	eventIndex int
	snapshot   proto.Message
}

type SnapshotStore struct {
	client    *dynamodb.Client
	table     string
	snapshots map[string]*snapshotEntry
}

func NewSnapshotStore(client *dynamodb.Client, table string) *SnapshotStore {
	return &SnapshotStore{
		client:    client,
		table:     table,
		snapshots: make(map[string]*snapshotEntry),
	}
}

func (s *SnapshotStore) GetSnapshot(actorName string) (snapshot interface{}, eventIndex int, ok bool) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.table),
		KeyConditionExpression: aws.String("actorName = :actorName"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":actorName": &types.AttributeValueMemberS{Value: actorName},
		},
		ScanIndexForward: aws.Bool(false), // 逆順にソート
		Limit:            aws.Int32(1),    // 最新の1レコードのみ取得
	}

	result, err := s.client.Query(context.Background(), input)
	if err != nil {
		return nil, 0, false
	}

	if len(result.Items) == 0 {
		return nil, 0, false
	}

	item := result.Items[0]

	var snapshotData map[string]interface{}
	err = attributevalue.UnmarshalMap(item, &snapshotData)
	if err != nil {
		return nil, 0, false
	}

	// Goでは、UnmarshalMapすると数値はfloat64になるが、取得できているかを型assertionで確認する
	eventIndexStr := fmt.Sprintf("%.0f", snapshotData["eventIndex"].(float64))
	eventIndex, err = strconv.Atoi(eventIndexStr)
	if err != nil {
		return nil, 0, false
	}

	snapshotBytes, ok := snapshotData["payload"].([]byte)
	if !ok {
		if snapshotData["payload"] == nil {
			// log something
			return nil, 0, false
		}
		return nil, 0, false
	}

	snapshot = &Snapshot{}
	err = proto.Unmarshal(snapshotBytes, snapshot.(*Snapshot))
	if err != nil {
		return nil, 0, false
	}

	return snapshot, eventIndex, true
}

func (s *SnapshotStore) PersistSnapshot(actorName string, snapshotIndex int, snapshot protoreflect.ProtoMessage) {
	// TODO: Implement persisting the snapshot to the database
	// 1. Convert the snapshot to a format suitable for storing in the database (e.g., JSON, binary)
	// 2. Store the converted snapshot in the database, associating it with the actorName and snapshotIndex
}

func (s *SnapshotStore) DeleteSnapshots(actorName string, inclusiveToIndex int) {
	// TODO: Implement deleting snapshots from the database
	// 1. Delete all snapshots for the given actorName up to and including the specified inclusiveToIndex
}
