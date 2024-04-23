package grain_test

import (
	"testing"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func TestAccountGrain_Add(t *testing.T) {
	// NewTestProviderではtestができなかった
	// やむなく、cluster packageを使う
	c := cluster.StartNode("my-cluster6332", 6332)
	defer c.Shutdown(true)

	// protoもやむなく使う
	proto.AccountFactory(func() proto.Account {
		return &grain.AccountGrain{}
	})

	accountGrainClient := proto.GetAccountGrainClient(c, "test_account")

	req := &proto.AccountRegisterRequest{Id: "mockID", Email: "test@account.test"}
	resp, err := accountGrainClient.Register(req)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.Id != "mockID" {
		t.Errorf("Expected count: mockID, but got: %s", resp.Id)
	}
}
