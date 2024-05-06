package rmu

import (
	"context"
	"fmt"

	dynamodbevents "github.com/aws/aws-lambda-go/events"
	"google.golang.org/protobuf/proto"
	_ "github.com/go-sql-driver/mysql"
	"log/slog"
	"strings"
	p "github.com/tkhrk1010/protoactor-go-persistence-dynamodb/persistence"
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
	slog.Info("start UpdateReadModel")

	for _, record := range event.Records {
		slog.Info("Processing request data for event GetId %s, type %s.", record.EventID, record.EventName)
		attributeValues := record.Change.NewImage
		payloadBytes := convertToBytes(attributeValues["payload"])
		typeValueStr, err := getTypeString(payloadBytes)
		if err != nil {
			return err
		}
		slog.Info(fmt.Sprintf("typeValueStr = %s", typeValueStr))
		if strings.HasPrefix(typeValueStr, "windSpeed") {
			event, err := convertWindSpeedEvent(payloadBytes)
			if err != nil {
				slog.Info(fmt.Sprintf("convertWindSpeedEvent: err = %v", err))
				return err
			}
			// debug
			slog.Info(fmt.Sprintf("event = %v, type = %T", event, event))

			switch event.(type) {
			case *WindSpeedCreated:
				//debug
				slog.Info(fmt.Sprintf("UpdateReadModel > WindSpeedCreated: %v", event))

				ev := event.(*WindSpeedCreated)
				err2 := createWindSpeed(ev, r)
				if err2 != nil {
					slog.Info(fmt.Sprintf("createWindSpeed: err = %v", err2))
					return err2
				}
			case *WindSpeedUpdated:
				ev := event.(*WindSpeedUpdated)
				err2 := updateWindSpeed(ev, r)
				if err2 != nil {
					return err2
				}
			default:
				return fmt.Errorf("unknown event type: %T", event)
			}
		}
		// Print new values for attributes of type String
		for name, value := range record.Change.NewImage {
			slog.Info(fmt.Sprintf("Attribute name: %s, value: %s", name, value.String()))
		}
	}
	return nil
}

func createWindSpeed(ev *WindSpeedCreated, r *ReadModelUpdater) error {
	slog.Info(fmt.Sprintf("createWindSpeed: start: ev = %v", ev))
	id := ev.GetId()
	value := ev.GetValue()
	occurredAt := ev.GetOccurredAt().AsTime()
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
	occurredAt := ev.GetOccurredAt().AsTime()
	err := r.dao.UpdateWindSpeed(id, value, occurredAt)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("updateWindSpeed: finished"))
	return nil
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
	message := &p.Event{}
	err := proto.Unmarshal(bytes, message)
	if err != nil {
		slog.Info(fmt.Sprintf("getTypeString: err = %v, %s", err, string(bytes)))
		return "", err
	}
	typeValue := message.GetType()
	if typeValue == "" {
			return "", fmt.Errorf("type is empty")
	}
	return typeValue, nil
}

func convertWindSpeedEvent(payloadBytes []byte) (Event, error) {
	message := &p.Event{}
	err := proto.Unmarshal(payloadBytes, message)
	if err != nil {
		slog.Info(fmt.Sprintf("convertWindSpeedEvent: proto unmarshal err = %v, %s", err, string(payloadBytes)))
		return nil, err
	}
	event, err := EventConverter(message)
	if err != nil {
		slog.Info(fmt.Sprintf("convertWindSpeedEvent: EventConverter err = %v, %s", err, string(payloadBytes)))
		return nil, err
	}
	return event, nil
}
