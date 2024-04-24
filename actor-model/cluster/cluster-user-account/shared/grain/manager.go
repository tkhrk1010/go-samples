package grain

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/domain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

type ManagerGrain struct {
	accountMap map[string]*actor.PID
}

type AccountActor struct {
	account domain.Account
}

func (a *AccountActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *proto.AccountIdResponse:
		a.account.ID = msg.Id
	case *proto.AccountResponse:
		a.account.ID = msg.Id
		a.account.Email = msg.Email
	case *proto.Noop:
		ctx.Respond(&proto.AccountResponse{Id: a.account.ID, Email: a.account.Email})
	}
}

func (m *ManagerGrain) ReceiveDefault(ctx cluster.GrainContext) {
}

func (m *ManagerGrain) Init(ctx cluster.GrainContext) {
	m.accountMap = make(map[string]*actor.PID)
}

func (m *ManagerGrain) Terminate(ctx cluster.GrainContext) {
}

func (m *ManagerGrain) CreateAccount(n *proto.CreateAccountRequest, ctx cluster.GrainContext) (*proto.AccountIdResponse, error) {
	// ManagerGrainのidentityがaccountのidになる
	id := ctx.Identity()
	account := domain.NewAccount(id, n.Email)
	accountActor := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor { return &AccountActor{account: *account} }))
	m.accountMap[account.ID] = accountActor
	return &proto.AccountIdResponse{Id: account.ID}, nil
}

func (m *ManagerGrain) GetAccount(n *proto.AccountIdResponse, ctx cluster.GrainContext) (*proto.AccountResponse, error) {
	accountActor := m.accountMap[n.Id]
	future := ctx.RequestFuture(accountActor, &proto.Noop{}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		return nil, err
	}
	return result.(*proto.AccountResponse), nil
}

func (m *ManagerGrain) GetAllAccountEmails(n *proto.Noop, ctx cluster.GrainContext) (*proto.EmailsResponse, error) {
	emails := make(map[string]string)
	for id, accountActor := range m.accountMap {
		future := ctx.RequestFuture(accountActor, &proto.Noop{}, 5*time.Second)
		result, err := future.Result()
		if err != nil {
			return nil, err
		}
		account := result.(*proto.AccountResponse)
		emails[id] = account.Email
	}
	return &proto.EmailsResponse{Emails: emails}, nil
}
