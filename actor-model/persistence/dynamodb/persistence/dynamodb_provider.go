package persistence

import (
	"github.com/asynkron/protoactor-go/persistence"
)


// ProviderState is an object containing the implementation for the provider
type ProviderState struct {
	snapshotStore persistence.SnapshotStore
	eventStore    persistence.EventStore
}

// NewProviderState creates a new instance of ProviderState
func NewProviderState(snapshotStore persistence.SnapshotStore, eventStore persistence.EventStore) *ProviderState {
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
	// TODO: Implement getting the current state
	return p
}

// Restart restarts the provider
func (p *ProviderState) Restart() {}

// GetSnapshotInterval returns the snapshot interval
func (p *ProviderState) GetSnapshotInterval() int {
	// TODO: Implement getting the snapshot interval
	return 3
}