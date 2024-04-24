package presentation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/presentation"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func TestRegisterAccount(t *testing.T) {
	// Start the cluster
	c := cluster.StartNode("my-cluster6333", 6333)
	defer c.Shutdown(true)

	// Register an account
	id := presentation.RegisterAccount(c, "email1@account.test")

	// Get the ManagerGrain client
	managerGrainClient := proto.GetManagerGrainClient(c, "singleManagerGrain")

	// Retrieve the account
	resp, err := managerGrainClient.GetAccount(&proto.AccountIdResponse{Id: id})

	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, id, resp.Id, "Account ID should match")
	assert.Equal(t, "email1@account.test", resp.Email, "Account email should match")
}

func TestGetAllAccounts(t *testing.T) {
	// Start the cluster
	c := cluster.StartNode("my-cluster6334", 6334)
	defer c.Shutdown(true)

	// Register multiple accounts
	id1 := presentation.RegisterAccount(c, "email1@account.test")
	id2 := presentation.RegisterAccount(c, "email2@account.test")
	id3 := presentation.RegisterAccount(c, "email3@account.test")

	// Get all accounts
	presentation.GetAllAccounts(c)

	// Get the ManagerGrain client
	managerGrainClient := proto.GetManagerGrainClient(c, "singleManagerGrain")

	// Retrieve the emails
	resp, err := managerGrainClient.GetAllAccountEmails(&proto.Noop{})

	assert.NoError(t, err, "Unexpected error")

	expected := map[string]string{
		id1: "email1@account.test",
		id2: "email2@account.test",
		id3: "email3@account.test",
	}

	assert.Equal(t, len(expected), len(resp.Emails), "Number of emails should match")

	for id, email := range expected {
		assert.Equal(t, email, resp.Emails[id], "Email should match")
	}
}