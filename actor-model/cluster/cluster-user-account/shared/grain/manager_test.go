package grain_test

import (
	"testing"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
	c "github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
)

func TestManagerGrain_RegisterGrain(t *testing.T) {
	c := c.StartNode(6330)
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
	c := c.StartNode(6330)
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
// func TestManagerGrain_BroadcastGetCounts(t *testing.T) {
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

// 	req := &proto.NumberRequest{Number: 10}
// 	_, err := accountGrainClient.Add(req)

// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	resp, err := managerGrainClient.BroadcastGetCounts(&proto.Noop{})

// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	expected := map[string]int64{"test_account": 10}
// 	if len(resp.Totals) != len(expected) {
// 		t.Errorf("Expected totals length: %d, but got: %d", len(expected), len(resp.Totals))
// 	}

// 	for grainAddress, total := range expected {
// 		if resp.Totals[grainAddress] != total {
// 			t.Errorf("Expected total for %s: %d, but got: %d", grainAddress, total, resp.Totals[grainAddress])
// 		}
// 	}
// }