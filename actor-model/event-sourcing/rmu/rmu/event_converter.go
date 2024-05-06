// 簡単のため、rmuに直置きする
package rmu

import (
	"fmt"
	"log/slog"
)

func EventConverter(m map[string]interface{}) (Event, error) {
	slog.Info(fmt.Sprintf("EventConverter: %v", m))
	eventId := m["id"].(string)
	occurredAt := uint64(m["occurred_at"].(float64))
	switch m["type"].(string) {
	case "windSpeedCollect":
		dataMap, ok := m["Data"].(map[string]interface{})
		if !ok {
			fmt.Println("Invalid data format")
			return nil, fmt.Errorf("invalid data format")
		}
		value, ok := dataMap["value"].(float64)
		if !ok {
			fmt.Println("Invalid value format")
			return nil, fmt.Errorf("invalid value format")
		}

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
		dataMap, ok := m["Data"].(map[string]interface{})
		if !ok {
			fmt.Println("Invalid data format")
			return nil, fmt.Errorf("invalid data format")
		}
		value, ok := dataMap["value"].(float64)
		if !ok {
			fmt.Println("Invalid value format")
			return nil, fmt.Errorf("invalid value format")
		}
		event := NewWindSpeedUpdatedFrom(
			eventId,
			value,
			occurredAt,
		)
		return &event, nil

	default:
		return nil, fmt.Errorf("unknown event type: %s", m["type"].(string))
	}
}
