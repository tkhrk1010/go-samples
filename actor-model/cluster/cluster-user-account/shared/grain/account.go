// virtual actorを提供する。accountを呼び出すためのactor
// accountがgrainの必要なさそう。actorで十分では？
// そうするためには、managerGrainがactorを管理する必要がある
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
	// register with the manager
	managerGrain := proto.GetManagerGrainClient(ctx.Cluster(), "singleManagerGrain")
	managerGrain.RegisterGrain(&proto.RegisterMessage{GrainId: ctx.Identity()})
}

func (s *AccountGrain) Terminate(ctx cluster.GrainContext) {
	// deregister with the manager
	managerGrain := proto.GetManagerGrainClient(ctx.Cluster(), "singleManagerGrain")
	managerGrain.DeregisterGrain(&proto.RegisterMessage{GrainId: ctx.Identity()})
}

func (s *AccountGrain) Register(n *proto.AccountRegisterRequest, ctx cluster.GrainContext) (*proto.AccountIdResponse, error) {
	a := domain.NewAccount(n.Id, n.Email)
	s.account = *a
	return &proto.AccountIdResponse{Id: s.account.ID}, nil
}

func (s *AccountGrain) GetCurrent(n *proto.Noop, ctx cluster.GrainContext) (*proto.AccountResponse, error) {
	return &proto.AccountResponse{Id: s.account.ID, Email: s.account.Email}, nil
}
