package rmu

import (
	"context"
	"encoding/json"
	"fmt"

	dynamodbevents "github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
	"log/slog"
	"strings"
	"time"
)

type ReadModelUpdater struct {
	dao WindSpeedDao
}

// NewReadModelUpdater is a constructor for ReadModelUpdater.
func NewReadModelUpdater(dao WindSpeedDao) ReadModelUpdater {
	return ReadModelUpdater{dao}
}

// TODO: 楽観ロックかかってないかも。dynamodb streamsから同じeventが複数回来た場合、2回目以降はエラーになるかも
// eventIndexを使って、eventが処理済みかどうかを判定する必要があるかも

// UpdateReadModel processes events from DynamoDB stream and updates the read model.
func (r *ReadModelUpdater) UpdateReadModel(ctx context.Context, event dynamodbevents.DynamoDBEvent) error {
	for _, record := range event.Records {
		slog.Info("Processing request data for event GetId %s, type %s.", record.EventID, record.EventName)
		attributeValues := record.Change.NewImage
		payloadBytes := convertToBytes(attributeValues["payload"])
		typeValueStr, err := getTypeString(payloadBytes)
		if err != nil {
			return err
		}
		slog.Debug(fmt.Sprintf("typeValueStr = %s", typeValueStr))
		if strings.HasPrefix(typeValueStr, "windSpeed") {
			event, err := convertWindSpeedEvent(payloadBytes)
			if err != nil {
				return err
			}
			switch event.(type) {
			case *WindSpeedCreated:
				ev := event.(*WindSpeedCreated)
				err2 := createWindSpeed(ev, r)
				if err2 != nil {
					return err2
				}
			case *WindSpeedUpdated:
				ev := event.(*WindSpeedUpdated)
				err2 := updateWindSpeed(ev, r)
				if err2 != nil {
					return err2
				}
			default:
			}
		}
		// Print new values for attributes of type String
		for name, value := range record.Change.NewImage {
			slog.Debug(fmt.Sprintf("Attribute name: %s, value: %s", name, value.String()))
		}
	}
	return nil
}

func createWindSpeed(ev *WindSpeedCreated, r *ReadModelUpdater) error {
	slog.Info(fmt.Sprintf("createWindSpeed: start: ev = %v", ev))
	id := ev.GetId()
	value := ev.GetValue()
	occurredAt := convertToTime(ev.GetOccurredAt())
	err := r.dao.InsertWindSpeed(id, value, occurredAt)
	if err != nil {
		return err
	}

	slog.Info("createWindSpeed: finished")
	return nil
}

func updateWindSpeed(ev *WindSpeedUpdated, r *ReadModelUpdater) error {
	slog.Info(fmt.Sprintf("updateWindSpeed: start: ev = %v", ev))
	id := ev.GetId()
	value := ev.GetValue()
	occurredAt := convertToTime(ev.GetOccurredAt())
	err := r.dao.UpdateWindSpeed(id, value, occurredAt)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("updateWindSpeed: finished"))
	return nil
}

func convertToTime(epoc uint64) time.Time {
	occurredAtUnix := int64(epoc) * int64(time.Millisecond)
	occurredAt := time.Unix(0, occurredAtUnix)
	return occurredAt
}

func convertToBytes(payloadAttr dynamodbevents.DynamoDBAttributeValue) []byte {
	var payloadBytes []byte
	if payloadAttr.DataType() == dynamodbevents.DataTypeBinary {
		payloadBytes = payloadAttr.Binary()
	} else if payloadAttr.DataType() == dynamodbevents.DataTypeString {
		payloadBytes = []byte(payloadAttr.String())
	}
	return payloadBytes
}

func getTypeString(bytes []byte) (string, error) {
	var parsed map[string]interface{}
	err := json.Unmarshal(bytes, &parsed)
	if err != nil {
		slog.Info(fmt.Sprintf("getTypeString: err = %v, %s", err, string(bytes)))
		return "", err
	}
	typeValue, ok := parsed["type"].(string)
	if !ok {
		return "", fmt.Errorf("type is not a string")
	}
	return typeValue, nil
}

func convertWindSpeedEvent(payloadBytes []byte) (Event, error) {
	var parsed map[string]interface{}
	err := json.Unmarshal(payloadBytes, &parsed)
	if err != nil {
		return nil, err
	}
	event, err := EventConverter(parsed)
	if err != nil {
		return nil, err
	}
	return event, nil
}
