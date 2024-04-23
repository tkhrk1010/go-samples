package presentation_test

import (
	"testing"
	"time"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/presentation"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func TestRegisterAccount(t *testing.T) {
	// Start the cluster
	c := cluster.StartNode(6330)
	defer c.Shutdown(true)

	// Register an account
	presentation.RegisterAccount(c, "test_account", 10)

	// Get the AccountGrain client
	accountGrainClient := proto.GetAccountGrainClient(c, "test_account")

	// Retrieve the current count
	resp, err := accountGrainClient.GetCurrent(&proto.Noop{})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.Number != 10 {
		t.Errorf("Expected count: 10, but got: %d", resp.Number)
	}
}

func TestGetAllAccounts(t *testing.T) {
	// Start the cluster
	c := cluster.StartNode(6330)
	defer c.Shutdown(true)

	// Register multiple accounts
	presentation.RegisterAccount(c, "account1", 10)
	presentation.RegisterAccount(c, "account2", 20)
	presentation.RegisterAccount(c, "account3", 30)

	// Wait for the cluster to stabilize
	time.Sleep(1 * time.Second)

	// Get all accounts
	presentation.GetAllAccounts(c)

	// Get the ManagerGrain client
	managerGrainClient := proto.GetManagerGrainClient(c, "singleManagerGrain")

	// Retrieve the totals
	resp, err := managerGrainClient.BroadcastGetCounts(&proto.Noop{})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := map[string]int64{
		"account1": 10,
		"account2": 20,
		"account3": 30,
	}

	if len(resp.Totals) != len(expected) {
		t.Errorf("Expected totals length: %d, but got: %d", len(expected), len(resp.Totals))
	}

	for grainId, count := range expected {
		if resp.Totals[grainId] != count {
			t.Errorf("Expected count for %s: %d, but got: %d", grainId, count, resp.Totals[grainId])
		}
	}
}