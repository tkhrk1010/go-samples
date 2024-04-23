package grain_test

import (
	"testing"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func TestManagerGrain_RegisterGrain(t *testing.T) {
	c := cluster.StartNode("my-cluster6333", 6333)
	defer c.Shutdown(true)

	proto.ManagerFactory(func() proto.Manager {
		return &grain.ManagerGrain{}
	})

	managerGrainClient := proto.GetManagerGrainClient(c, "test_manager")

	req := &proto.RegisterMessage{GrainId: "test_account"}
	_, err := managerGrainClient.RegisterGrain(req)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestManagerGrain_DeregisterGrain(t *testing.T) {
	c := cluster.StartNode("my-cluster6334", 6334)
	defer c.Shutdown(true)

	proto.ManagerFactory(func() proto.Manager {
		return &grain.ManagerGrain{}
	})

	managerGrainClient := proto.GetManagerGrainClient(c, "test_manager")

	req := &proto.RegisterMessage{GrainId: "test_account"}
	_, err := managerGrainClient.DeregisterGrain(req)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// // FIXME: 登録か取得に失敗している。
// func TestManagerGrain_GetAllAccountEmails(t *testing.T) {
// 	c := c.StartNode(6330)
// 	defer c.Shutdown(true)

// 	proto.ManagerFactory(func() proto.Manager {
// 		return &grain.ManagerGrain{}
// 	})

// 	managerGrainClient := proto.GetManagerGrainClient(c, "test_manager")

// 	proto.AccountFactory(func() proto.Account {
// 		return &grain.AccountGrain{}
// 	})

// 	accountGrainClient := proto.GetAccountGrainClient(c, "test_account")

// 	req := &proto.AccountRegisterRequest{Id: 10}
// 	_, err := accountGrainClient.Add(req)

// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	resp, err := managerGrainClient.GetAllAccountEmails(&proto.Noop{})

// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	expected := map[string]int64{"test_account": 10}
// 	if len(resp.Emails) != len(expected) {
// 		t.Errorf("Expected emails length: %d, but got: %d", len(expected), len(resp.Emails))
// 	}

// 	for grainAddress, total := range expected {
// 		if resp.Emails[grainAddress] != total {
// 			t.Errorf("Expected total for %s: %d, but got: %d", grainAddress, total, resp.Emails[grainAddress])
// 		}
// 	}
// }
