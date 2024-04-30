package persistence_test

import (
	"testing"

	"github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/persistence"
)

func TestNewProviderState(t *testing.T) {
	ps := persistence.NewProviderState()
	if ps == nil {
		t.Error("NewProviderState returned nil")
	}
}

func TestGetState(t *testing.T) {
	ps := persistence.NewProviderState()
	state := ps.GetState()
	if state != ps {
		t.Error("GetState should return the same instance")
	}
}

func TestRestart(t *testing.T) {
	ps := persistence.NewProviderState()
	ps.Restart()
	// TODO: Add assertions for the expected behavior after restarting
}

func TestGetSnapshotInterval(t *testing.T) {
	ps := persistence.NewProviderState()
	interval := ps.GetSnapshotInterval()
	if interval != 0 {
		t.Errorf("GetSnapshotInterval should return 0, got: %d", interval)
	}
}