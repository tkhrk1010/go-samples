package persistence

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type EventStore struct {
	// TODO: Add necessary fields, such as database connection or configuration
}

func NewEventStore(/* TODO: Add necessary parameters */) *EventStore {
	return &EventStore{
		// TODO: Initialize fields
	}
}

func (e *EventStore) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{})) {
	// TODO: Implement getting events from the database
	// 1. Query the database to retrieve events for the given actorName within the specified eventIndexStart and eventIndexEnd range
	// 2. For each event found, call the provided callback function with the event
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