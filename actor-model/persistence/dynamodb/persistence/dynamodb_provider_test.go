package persistence_test

import (
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
	p "github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/persistence"
)

type mockSnapshotStore struct{}

func (m *mockSnapshotStore) GetSnapshot(actorName string) (snapshot interface{}, eventIndex int, ok bool) {
	return nil, 0, false
}

func (m *mockSnapshotStore) PersistSnapshot(actorName string, snapshotIndex int, snapshot protoreflect.ProtoMessage) {
}

func (m *mockSnapshotStore) DeleteSnapshots(actorName string, inclusiveToIndex int) {
}

type mockEventStore struct{}

func (m *mockEventStore) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{})) {
}

func (m *mockEventStore) PersistEvent(actorName string, eventIndex int, event protoreflect.ProtoMessage) {
}

func (m *mockEventStore) DeleteEvents(actorName string, inclusiveToIndex int) {
}

func TestNewProviderState(t *testing.T) {
	snapshotStore := &mockSnapshotStore{}
	eventStore := &mockEventStore{}
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
	snapshotStore := &mockSnapshotStore{}
	eventStore := &mockEventStore{}
	ps := p.NewProviderState(snapshotStore, eventStore)
	state := ps.GetState()
	if state != ps {
		t.Error("GetState should return the same instance")
	}
}

func TestRestart(t *testing.T) {
	snapshotStore := &mockSnapshotStore{}
	eventStore := &mockEventStore{}
	ps := p.NewProviderState(snapshotStore, eventStore)
	ps.Restart()
}

func TestGetSnapshotInterval(t *testing.T) {
	snapshotStore := &mockSnapshotStore{}
	eventStore := &mockEventStore{}
	ps := p.NewProviderState(snapshotStore, eventStore)
	interval := ps.GetSnapshotInterval()
	if interval != 3 {
		t.Errorf("GetSnapshotInterval should return 100, got: %d", interval)
	}
}