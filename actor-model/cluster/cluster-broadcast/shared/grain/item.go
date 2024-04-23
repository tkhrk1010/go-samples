// virtual actorを提供する。itemを呼び出すためのactor
// itemというよりitem(買い物している人)かな
package grain

import (
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/domain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/proto"
)

type ItemGrain struct {
	item domain.Item
}

func (s *ItemGrain) ReceiveDefault(ctx cluster.GrainContext) {
}

func (s *ItemGrain) Init(ctx cluster.GrainContext) {
	s.item.Count = 0

	// register with the cart
	cartGrain := proto.GetCartGrainClient(ctx.Cluster(), "singleCartGrain")
	cartGrain.RegisterGrain(&proto.RegisterMessage{GrainId: ctx.Identity()})
}

func (s *ItemGrain) Terminate(ctx cluster.GrainContext) {
	// deregister with the cart
	cartGrain := proto.GetCartGrainClient(ctx.Cluster(), "singleCartGrain")
	cartGrain.DeregisterGrain(&proto.RegisterMessage{GrainId: ctx.Identity()})
}

func (s *ItemGrain) Add(n *proto.NumberRequest, ctx cluster.GrainContext) (*proto.CountResponse, error) {
	s.item.Add(n.Number)
	return &proto.CountResponse{Number: s.item.Count}, nil
}

func (s *ItemGrain) Remove(n *proto.NumberRequest, ctx cluster.GrainContext) (*proto.CountResponse, error) {
	s.item.Remove(n.Number)
	return &proto.CountResponse{Number: s.item.Count}, nil
}

func (s *ItemGrain) GetCurrent(n *proto.Noop, ctx cluster.GrainContext) (*proto.CountResponse, error) {
	return &proto.CountResponse{Number: s.item.Count}, nil
}
