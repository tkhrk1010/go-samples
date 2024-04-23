package grain

import (
	"fmt"
	"strings"

	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/proto"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/domain"
)

type CartGrain struct {
	cart domain.Cart
}

func (t *CartGrain) ReceiveDefault(ctx cluster.GrainContext) {
}

func (t *CartGrain) Init(ctx cluster.GrainContext) {
	t.cart.Init()
}

func (t *CartGrain) Terminate(ctx cluster.GrainContext) {
}

func (t *CartGrain) RegisterGrain(n *proto.RegisterMessage, ctx cluster.GrainContext) (*proto.Noop, error) {
	parts := strings.Split(n.GrainId, "/")
	grainID := parts[len(parts)-1]
	t.cart.RegisterItem(grainID)
	return &proto.Noop{}, nil
}

func (t *CartGrain) DeregisterGrain(n *proto.RegisterMessage, ctx cluster.GrainContext) (*proto.Noop, error) {
	t.cart.DeregisterItem(n.GrainId)
	return &proto.Noop{}, nil
}

func (t *CartGrain) BroadcastGetCounts(n *proto.Noop, ctx cluster.GrainContext) (*proto.TotalsResponse, error) {
	totals := map[string]int64{}
	// item nameをgrainAddressとして使用している
	for grainAddress := range t.cart.ItemMap {
		itemGrain := proto.GetItemGrainClient(ctx.Cluster(), grainAddress)
		grainTotal, err := itemGrain.GetCurrent(&proto.Noop{})
		if err != nil {
			fmt.Sprintf("Grain %s issued an error : %s", grainAddress, err)
		}
		fmt.Sprintf("Grain %s - %v", grainAddress, grainTotal.Number)
		totals[grainAddress] = grainTotal.Number
	}

	return &proto.TotalsResponse{Totals: totals}, nil
}
