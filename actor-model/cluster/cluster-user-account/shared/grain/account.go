// virtual actorを提供する。accountを呼び出すためのactor
// accountというよりaccount(買い物している人)かな
package grain

import (
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/domain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

type AccountGrain struct {
	account domain.Account
}

func (s *AccountGrain) ReceiveDefault(ctx cluster.GrainContext) {
}

func (s *AccountGrain) Init(ctx cluster.GrainContext) {
	s.account.Count = 0

	// register with the manager
	managerGrain := proto.GetManagerGrainClient(ctx.Cluster(), "singleManagerGrain")
	managerGrain.RegisterGrain(&proto.RegisterMessage{GrainId: ctx.Identity()})
}

func (s *AccountGrain) Terminate(ctx cluster.GrainContext) {
	// deregister with the manager
	managerGrain := proto.GetManagerGrainClient(ctx.Cluster(), "singleManagerGrain")
	managerGrain.DeregisterGrain(&proto.RegisterMessage{GrainId: ctx.Identity()})
}

func (s *AccountGrain) Add(n *proto.NumberRequest, ctx cluster.GrainContext) (*proto.CountResponse, error) {
	s.account.Add(n.Number)
	return &proto.CountResponse{Number: s.account.Count}, nil
}

func (s *AccountGrain) Remove(n *proto.NumberRequest, ctx cluster.GrainContext) (*proto.CountResponse, error) {
	s.account.Remove(n.Number)
	return &proto.CountResponse{Number: s.account.Count}, nil
}

func (s *AccountGrain) GetCurrent(n *proto.Noop, ctx cluster.GrainContext) (*proto.CountResponse, error) {
	return &proto.CountResponse{Number: s.account.Count}, nil
}
