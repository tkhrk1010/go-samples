package presentation_test

import (
	"testing"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/presentation"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func TestRegisterAccount(t *testing.T) {
	// Start the cluster
	c := cluster.StartNode("my-cluster6335", 6335)
	defer c.Shutdown(true)

	// Register an account
	id := presentation.RegisterAccount(c, "email1@account.test")

	// Get the AccountGrain client
	accountGrainClient := proto.GetAccountGrainClient(c, id)

	// Retrieve the current count
	resp, err := accountGrainClient.GetCurrent(&proto.Noop{})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.Id != id {
		t.Errorf("Expected count: %s, but got: %s", id, resp.Id)
	}
}

func TestGetAllAccounts(t *testing.T) {
	// Start the cluster
	c := cluster.StartNode("my-cluster6336", 6336)
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

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := map[string]string{
		id1: "email1@account.test",
		id2: "email2@account.test",
		id3: "email3@account.test",
	}

	if len(resp.Emails) != len(expected) {
		t.Errorf("Expected emails length: %d, but got: %d", len(expected), len(resp.Emails))
	}

	for grainId, count := range expected {
		if resp.Emails[grainId] != count {
			t.Errorf("Expected count for %s: %s, but got: %s", grainId, count, resp.Emails[grainId])
		}
	}
}
