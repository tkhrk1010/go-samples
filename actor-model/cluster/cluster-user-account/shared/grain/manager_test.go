package grain_test

import (
	"testing"

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
	c := cluster.StartNode("my-cluster6332", 6332)
	defer c.Shutdown(true)
	managerGrain := proto.GetManagerGrainClient(c, "test-grain")

	// Test CreateAccount
	createAccountResp, err := managerGrain.CreateAccount(&proto.CreateAccountRequest{Email: "testemail"})

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
