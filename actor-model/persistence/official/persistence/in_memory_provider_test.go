package provider

import (
	"testing"

	"github.com/asynkron/protoactor-go/persistence"
	"github.com/stretchr/testify/assert"
)

func TestNewProvider(t *testing.T) {
	snapshotInterval := 10

	provider := NewProvider(snapshotInterval)

	assert.NotNil(t, provider)
	assert.NotNil(t, provider.providerState)
	assert.Equal(t, snapshotInterval, provider.providerState.GetSnapshotInterval())
}

func TestProvider_GetState(t *testing.T) {
	snapshotInterval := 10

	provider := &Provider{
		providerState: persistence.NewInMemoryProvider(snapshotInterval),
	}

	state := provider.GetState()

	assert.NotNil(t, state)
	assert.Equal(t, snapshotInterval, state.GetSnapshotInterval())
}