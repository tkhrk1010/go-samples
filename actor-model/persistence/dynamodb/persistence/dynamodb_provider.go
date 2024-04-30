package persistence

import (
	"github.com/asynkron/protoactor-go/persistence"
	"google.golang.org/protobuf/reflect/protoreflect"
)


// ProviderState is an object containing the implementation for the provider
type ProviderState struct {
	snapshotStore persistence.SnapshotStore
	eventStore    persistence.EventStore
}

// NewProviderState creates a new instance of ProviderState
func NewProviderState(snapshotStore *SnapshotStore, eventStore *EventStore) *ProviderState {
	return &ProviderState{
		snapshotStore: snapshotStore,
		eventStore:    eventStore,
	}
}

func (p *ProviderState) GetSnapshotStore() persistence.SnapshotStore {
	return p.snapshotStore
}

func (p *ProviderState) GetEventStore() persistence.EventStore {
	return p.eventStore
}

// GetState returns the current state of the provider
func (p *ProviderState) GetState() *ProviderState {
	return p
}

// Restart restarts the provider
func (p *ProviderState) Restart() {}

// GetSnapshotInterval returns the snapshot interval
func (p *ProviderState) GetSnapshotInterval() int {
	// TODO: Implement getting the snapshot interval
	return 3
}

func (p *ProviderState) GetSnapshot(actorName string) (snapshot interface{}, eventIndex int, ok bool) {
	return p.snapshotStore.GetSnapshot(actorName)
}

func (p *ProviderState) PersistSnapshot(actorName string, snapshotIndex int, snapshot protoreflect.ProtoMessage) {
	p.snapshotStore.PersistSnapshot(actorName, snapshotIndex, snapshot)
}

func (p *ProviderState) DeleteSnapshots(actorName string, inclusiveToIndex int) {
	p.snapshotStore.DeleteSnapshots(actorName, inclusiveToIndex)
}

func (p *ProviderState) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{})) {
	p.eventStore.GetEvents(actorName, eventIndexStart, eventIndexEnd, callback)
}

func (p *ProviderState) PersistEvent(actorName string, eventIndex int, event protoreflect.ProtoMessage) {
	p.eventStore.PersistEvent(actorName, eventIndex, event)
}

func (p *ProviderState) DeleteEvents(actorName string, inclusiveToIndex int) {
	p.eventStore.DeleteEvents(actorName, inclusiveToIndex)
}