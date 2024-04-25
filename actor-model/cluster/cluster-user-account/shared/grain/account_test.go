package grain_test

import (
	"testing"
	"github.com/stretchr/testify/assert"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func TestAccountGrain(t *testing.T) {
	// Setup
	proto.AccountFactory(func() proto.Account {
		return &grain.AccountGrain{}
	})
	c := cluster.StartNode("my-cluster6332", 6332)
	defer c.Shutdown(true)
	AmccountGrain := proto.GetAccountGrainClient(c, "test-grain")

	// Test CreateAccount
	createAccountResp, err := AmccountGrain.CreateAccount(&proto.CreateAccountRequest{Email: "testemail"})

	assert.NoError(t, err)
	assert.NotEmpty(t, createAccountResp.Id)

	// Test GetAccount
	getAccountResp, err := AmccountGrain.GetAccount(&proto.AccountIdResponse{Id: createAccountResp.Id})
	assert.NoError(t, err)
	assert.Equal(t, createAccountResp.Id, getAccountResp.Id)
	assert.Equal(t, "testemail", getAccountResp.Email)

	// Test GetAccountEmail
	res, err := AmccountGrain.GetAccountEmail(&proto.Noop{})
	assert.NoError(t, err)
	assert.Equal(t, "testemail", res.Email)
}
