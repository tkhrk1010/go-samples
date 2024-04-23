package grain

import (
	"fmt"
	"strings"

	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/domain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

type ManagerGrain struct {
	manager domain.Manager
}

func (t *ManagerGrain) ReceiveDefault(ctx cluster.GrainContext) {
}

func (t *ManagerGrain) Init(ctx cluster.GrainContext) {
	t.manager.Init()
}

func (t *ManagerGrain) Terminate(ctx cluster.GrainContext) {
}

func (t *ManagerGrain) RegisterGrain(n *proto.RegisterMessage, ctx cluster.GrainContext) (*proto.Noop, error) {
	parts := strings.Split(n.GrainId, "/")
	grainID := parts[len(parts)-1]
	t.manager.RegisterAccount(grainID)
	return &proto.Noop{}, nil
}

func (t *ManagerGrain) DeregisterGrain(n *proto.RegisterMessage, ctx cluster.GrainContext) (*proto.Noop, error) {
	t.manager.DeregisterAccount(n.GrainId)
	return &proto.Noop{}, nil
}

func (t *ManagerGrain) BroadcastGetCounts(n *proto.Noop, ctx cluster.GrainContext) (*proto.TotalsResponse, error) {
	totals := map[string]int64{}
	// account nameをgrainAddressとして使用している
	for grainAddress := range t.manager.AccountMap {
		accountGrain := proto.GetAccountGrainClient(ctx.Cluster(), grainAddress)
		grainTotal, err := accountGrain.GetCurrent(&proto.Noop{})
		if err != nil {
			fmt.Sprintf("Grain %s issued an error : %s", grainAddress, err)
		}
		fmt.Sprintf("Grain %s - %v", grainAddress, grainTotal.Number)
		totals[grainAddress] = grainTotal.Number
	}

	return &proto.TotalsResponse{Totals: totals}, nil
}
