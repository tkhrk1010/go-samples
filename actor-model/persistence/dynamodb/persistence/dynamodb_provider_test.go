package persistence_test

import (
	"testing"

	p "github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/persistence"
)

func TestNewProviderState(t *testing.T) {
	snapshotStore := p.NewSnapshotStore()
	eventStore := p.NewEventStore()
	ps := p.NewProviderState(snapshotStore, eventStore)
	if ps == nil {
		t.Error("NewProviderState returned nil")
	}
	if ps.GetSnapshotStore() != snapshotStore {
		t.Error("ProviderState's snapshotStore is not set correctly")
	}
	if ps.GetEventStore() != eventStore {
		t.Error("ProviderState's eventStore is not set correctly")
	}
}

func TestGetState(t *testing.T) {
	snapshotStore := p.NewSnapshotStore()
	eventStore := p.NewEventStore()
	ps := p.NewProviderState(snapshotStore, eventStore)
	state := ps.GetState()
	if state != ps {
		t.Error("GetState should return the same instance")
	}
}

func TestRestart(t *testing.T) {
	snapshotStore := p.NewSnapshotStore()
	eventStore := p.NewEventStore()
	ps := p.NewProviderState(snapshotStore, eventStore)
	ps.Restart()
}

func TestGetSnapshotInterval(t *testing.T) {
	snapshotStore := p.NewSnapshotStore()
	eventStore := p.NewEventStore()
	ps := p.NewProviderState(snapshotStore, eventStore)
	interval := ps.GetSnapshotInterval()
	if interval != 3 {
		t.Errorf("GetSnapshotInterval should return 100, got: %d", interval)
	}
}