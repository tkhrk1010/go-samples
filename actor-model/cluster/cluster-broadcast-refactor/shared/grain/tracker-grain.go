package grain

import (
	"fmt"
	"strings"

	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/proto"
)

type TrackGrain struct {
	grainsMap map[string]bool
}

func (t *TrackGrain) ReceiveDefault(ctx cluster.GrainContext) {
}

func (t *TrackGrain) Init(ctx cluster.GrainContext) {
	t.grainsMap = map[string]bool{}
}

func (t *TrackGrain) Terminate(ctx cluster.GrainContext) {
}

func (t *TrackGrain) RegisterGrain(n *proto.RegisterMessage, ctx cluster.GrainContext) (*proto.Noop, error) {
	parts := strings.Split(n.GrainId, "/")
	grainID := parts[len(parts)-1]
	t.grainsMap[grainID] = true
	return &proto.Noop{}, nil
}

func (t *TrackGrain) DeregisterGrain(n *proto.RegisterMessage, ctx cluster.GrainContext) (*proto.Noop, error) {
	delete(t.grainsMap, n.GrainId)
	return &proto.Noop{}, nil
}

func (t *TrackGrain) BroadcastGetCounts(n *proto.Noop, ctx cluster.GrainContext) (*proto.TotalsResponse, error) {
	totals := map[string]int64{}
	for grainAddress := range t.grainsMap {
		calcGrain := proto.GetCalculatorGrainClient(ctx.Cluster(), grainAddress)
		grainTotal, err := calcGrain.GetCurrent(&proto.Noop{})
		if err != nil {
			fmt.Sprintf("Grain %s issued an error : %s", grainAddress, err)
		}
		fmt.Sprintf("Grain %s - %v", grainAddress, grainTotal.Number)
		totals[grainAddress] = grainTotal.Number
	}

	return &proto.TotalsResponse{Totals: totals}, nil
}
