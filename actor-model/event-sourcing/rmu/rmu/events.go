// 本来は、domain/events.goにWindSpeedCreated, WindSpeedUpdatedなどのイベントを定義しているのがわかりやすいが、ここでは簡単に同じpackageで済ます
package rmu

import ()

type Event interface{}

type WindSpeedCreated struct {
	Event
	id        string
	value     float64
	createdAt uint64
}

type WindSpeedUpdated struct {
	Event
	id        string
	value     float64
	updatedAt uint64
}

func NewWindSpeedCreatedFrom(id string, value float64, createdAt uint64) WindSpeedCreated {
	return WindSpeedCreated{
		id:        id,
		value:     value,
		createdAt: createdAt,
	}
}

func NewWindSpeedUpdatedFrom(id string, value float64, updatedAt uint64) WindSpeedUpdated {
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

func (e *WindSpeedCreated) GetOccurredAt() uint64 {
	return e.createdAt
}

func (e *WindSpeedUpdated) GetId() string {
	return e.id
}

func (e *WindSpeedUpdated) GetValue() float64 {
	return e.value
}

func (e *WindSpeedUpdated) GetOccurredAt() uint64 {
	return e.updatedAt
}
