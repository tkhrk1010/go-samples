package grain_test

import (
	"testing"
	"time"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	// "github.com/asynkron/protoactor-go/cluster"
	// "github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	// "github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	// "github.com/asynkron/protoactor-go/remote"
	"github.com/stretchr/testify/assert"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
)

func TestManagerGrain(t *testing.T) {
	// Setup
	proto.ManagerFactory(func() proto.Manager {
		return &grain.ManagerGrain{}
	})
	c := cluster.StartNode("test-cluster", 6330)
	defer c.Shutdown(true)
	managerGrain := proto.GetManagerGrainClient(c, "test-grain")

	// Test CreateAccount
	createAccountResp, err := managerGrain.CreateAccount(&proto.Noop{})

	assert.NoError(t, err)
	assert.NotEmpty(t, createAccountResp.Id)

	// Test GetAccount
	getAccountResp, err := managerGrain.GetAccount(&proto.AccountIdResponse{Id: createAccountResp.Id})
	assert.NoError(t, err)
	assert.Equal(t, createAccountResp.Id, getAccountResp.Id)
	assert.Equal(t, "testemail", getAccountResp.Email)

	// Test GetAllAccountEmails
	getAllEmailsResp, err := managerGrain.GetAllAccountEmails(&proto.Noop{})
	assert.NoError(t, err)
	assert.Len(t, getAllEmailsResp.Emails, 1)
	assert.Equal(t, "testemail", getAllEmailsResp.Emails[createAccountResp.Id])
}

func TestAccountActor(t *testing.T) {
	// Setup
	system := actor.NewActorSystem()
	rootContext := system.Root
	props := actor.PropsFromProducer(func() actor.Actor { return &grain.AccountActor{} })
	pid := system.Root.Spawn(props)

	// Test Receive with AccountIdResponse
	rootContext.Send(pid, &proto.AccountIdResponse{Id: "test-id"})

	// Test Receive with AccountResponse
	rootContext.Send(pid, &proto.AccountResponse{Id: "test-id", Email: "test-email"})

	// Test Receive with Noop
	future := system.Root.RequestFuture(pid, &proto.Noop{}, 5*time.Second)
	response, err := future.Result()
	assert.NoError(t, err)
	assert.Equal(t, &proto.AccountResponse{Id: "test-id", Email: "test-email"}, response)

	// Cleanup
	system.Shutdown()
}
