// 本来は、domain/events.goにWindSpeedCreated, WindSpeedUpdatedなどのイベントを定義しているのがわかりやすいが、ここでは簡単に同じpackageで済ます
package rmu

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Event interface{}

type WindSpeedCreated struct {
	Event
	id        string
	value     float64
	createdAt *timestamppb.Timestamp
}

type WindSpeedUpdated struct {
	Event
	id        string
	value     float64
	updatedAt *timestamppb.Timestamp
}

func NewWindSpeedCreatedFrom(id string, value float64, createdAt *timestamppb.Timestamp) WindSpeedCreated {
	return WindSpeedCreated{
		id:        id,
		value:     value,
		createdAt: createdAt,
	}
}

func NewWindSpeedUpdatedFrom(id string, value float64, updatedAt *timestamppb.Timestamp) WindSpeedUpdated {
	return WindSpeedUpdated{
		id:        id,
		value:     value,
		updatedAt: updatedAt,
	}
}

func (e *WindSpeedCreated) GetId() string {
	return e.id
}

func (e *WindSpeedCreated) GetValue() float64 {
	return e.value
}

func (e *WindSpeedCreated) GetOccurredAt() *timestamppb.Timestamp {
	return e.createdAt
}

func (e *WindSpeedUpdated) GetId() string {
	return e.id
}

func (e *WindSpeedUpdated) GetValue() float64 {
	return e.value
}

func (e *WindSpeedUpdated) GetOccurredAt() *timestamppb.Timestamp {
	return e.updatedAt
}
