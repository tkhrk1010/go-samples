package cluster_test

import (
	"testing"
	"time"

	c "github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func TestStartNode(t *testing.T) {
	// Start the first node
	cluster1 := c.StartNode(6330)
	defer cluster1.Shutdown(true)

	// Start the second node
	cluster2 := c.StartNode(6331)
	defer cluster2.Shutdown(true)

	// Wait for the clusters to stabilize
	time.Sleep(3 * time.Second)

	// Get the member list from the first cluster
	members1 := cluster1.MemberList

	// Check if both nodes are present in the member list
	if members1.Length() != 2 {
		t.Errorf("Expected 2 members, but got: %d", members1.Length())
	}

	// Get the member list from the second cluster
	members2 := cluster2.MemberList

	// Check if both nodes are present in the member list
	if members2.Length() != 2 {
		t.Errorf("Expected 2 members, but got: %d", members2.Length())
	}

	// Get the AccountGrain client from the first cluster
	accountGrainClient := proto.GetAccountGrainClient(cluster1, "test_account")

	// Send a request to the AccountGrain
	req := &proto.NumberRequest{Number: 10}
	resp, err := accountGrainClient.Add(req)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.Number != 10 {
		t.Errorf("Expected count: 10, but got: %d", resp.Number)
	}

	// Get the ManagerGrain client from the second cluster
	managerGrainClient := proto.GetManagerGrainClient(cluster2, "test_manager")

	// Send a request to the ManagerGrain
	registerReq := &proto.RegisterMessage{GrainId: "test_account"}
	_, err = managerGrainClient.RegisterGrain(registerReq)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}