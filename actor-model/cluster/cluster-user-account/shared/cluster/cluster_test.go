package cluster_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	c "github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func TestStartNode(t *testing.T) {
	// Start the first node
	cluster1 := c.StartNode("my-cluster", 6330)
	defer cluster1.Shutdown(true)

	// Start the second node
	cluster2 := c.StartNode("my-cluster", 6331)
	defer cluster2.Shutdown(true)

	// Wait for the clusters to stabilize
	time.Sleep(5 * time.Second)

	// Get the member list from the first cluster
	members1 := cluster1.MemberList

	// Check if both nodes are present in the member list
	assert.Equal(t, 2, members1.Length(), "Expected 2 members in cluster1")

	// Get the member list from the second cluster
	members2 := cluster2.MemberList

	// Check if both nodes are present in the member list
	assert.Equal(t, 2, members2.Length(), "Expected 2 members in cluster2")

	// Check if the Account kind is registered in the first cluster
	AmccountKind := cluster1.GetClusterKind("Account")
	assert.NotNil(t, AmccountKind, "Account kind should be registered in cluster1")

	// Check if the Account kind is registered in the second cluster
	AmccountKind = cluster2.GetClusterKind("Account")
	assert.NotNil(t, AmccountKind, "Account kind should be registered in cluster2")

	// Create a test message
	testMsg := &proto.Noop{}

	// Send the test message to a non-existing grain in the first cluster
	_, err := cluster1.Request("non_existing_grain", "Account", testMsg)
	assert.Error(t, err, "Requesting a non-existing grain should return an error")

	// Send the test message to a non-existing grain in the second cluster
	_, err = cluster2.Request("non_existing_grain", "Account", testMsg)
	assert.Error(t, err, "Requesting a non-existing grain should return an error")
}
