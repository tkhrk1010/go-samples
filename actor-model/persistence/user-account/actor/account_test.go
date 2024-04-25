package actor_test

import (
	"testing"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/stretchr/testify/assert"
	a "github.com/tkhrk1010/go-samples/actor-model/persistence/user-account/actor"
)

// Inmemory persistence provider
type Provider struct {
	providerState persistence.ProviderState
}

func NewProvider(snapshotInterval int) *Provider {
	return &Provider{
		providerState: persistence.NewInMemoryProvider(snapshotInterval),
	}
}

func (p *Provider) GetState() persistence.ProviderState {
	return p.providerState
}

func TestActorReceive(t *testing.T) {
	// Create a test actor system and root context
	system := actor.NewActorSystem()
	rootContext := system.Root

	// Create a test provider
	provider := NewProvider(1)

	// Create a props with the test provider
	props := actor.PropsFromProducer(func() actor.Actor { return &a.AccountActor{} },
		actor.WithReceiverMiddleware(persistence.Using(provider)))

	// Spawn the actor
	pid, err := rootContext.SpawnNamed(props, "test-actor")
	assert.NoError(t, err)

	// 以下、panicになっていないことだけ確認する
	// Test case 1: Receive a Message
	rootContext.Send(pid, a.NewSignUpEvent("email1"))

	// Test case 2: Receive a RequestSnapshot
	rootContext.Send(pid, &persistence.RequestSnapshot{})

	// Test case 3: Receive a Snapshot
	rootContext.Send(pid, a.NewSignUpSnapshot("email2"))

	// Test case 4: Receive a ReplayComplete
	rootContext.Send(pid, &persistence.ReplayComplete{})
}
