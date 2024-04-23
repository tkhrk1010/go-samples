package grain_test

import (
	"testing"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
	c "github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
)

func TestAccountGrain_Add(t *testing.T) {
	// NewTestProviderではtestができなかった
	// やむなく、cluster packageを使う
	c := c.StartNode(6330)
	defer c.Shutdown(true)

	// protoもやむなく使う
	proto.AccountFactory(func() proto.Account {
		return &grain.AccountGrain{}
	})

	accountGrainClient := proto.GetAccountGrainClient(c, "test_account")

	req := &proto.NumberRequest{Number: 10}
	resp, err := accountGrainClient.Add(req)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.Number != 10 {
		t.Errorf("Expected count: 10, but got: %d", resp.Number)
	}
}