// 簡単のため、rmuに直置きする
package rmu

import (
	"encoding/json"
	"fmt"
	"log/slog"

	p "github.com/tkhrk1010/protoactor-go-persistence-dynamodb/persistence"
)

func EventConverter(m *p.Event) (Event, error) {
	slog.Info(fmt.Sprintf("EventConverter: %v", m))
	eventId := m.GetId()
	occurredAt := m.GetOccurredAt()
	eventType := m.GetType()
	// debug
	slog.Info(fmt.Sprintf("EventConverter event params: eventId: %v, type: %v, occurredAt: %v", eventId, eventType, occurredAt))
	switch eventType {
	case "windSpeedCollect":
		dataStr := m.GetData()
		// JSON 形式の文字列を構造体にアンマーシャリング
		var data struct {
			Value float64 `json:"value"`
		}
		err := json.Unmarshal([]byte(dataStr), &data)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, fmt.Errorf("error: %v", err)
		}
		value := data.Value

		// debug
		slog.Info(fmt.Sprintf("EventConverter: eventId: %v, value: %v, occurredAt: %v", eventId, value, occurredAt))

		event := NewWindSpeedCreatedFrom(
			eventId,
			value,
			occurredAt,
		)

		// debug
		slog.Info(fmt.Sprintf("EventConverter: windSpeedCollect: %v", event))

		return &event, nil

	case "windSpeedUpdate":
		dataStr := m.GetData()
		// JSON 形式の文字列を構造体にアンマーシャリング
		var data struct {
			Value float64 `json:"value"`
		}
		err := json.Unmarshal([]byte(dataStr), &data)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, fmt.Errorf("error: %v", err)
		}
		value := data.Value

		event := NewWindSpeedUpdatedFrom(
			eventId,
			value,
			occurredAt,
		)
		return &event, nil

	default:
		return nil, fmt.Errorf("unknown event type: %s", m.GetType())
	}
}
