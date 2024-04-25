package grain

import (
	"fmt"

	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/domain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

type ManagerGrain struct {
	accountMap map[string]*domain.Account
}

func (m *ManagerGrain) ReceiveDefault(ctx cluster.GrainContext) {
}

func (m *ManagerGrain) Init(ctx cluster.GrainContext) {
	m.accountMap = make(map[string]*domain.Account)
}

func (m *ManagerGrain) Terminate(ctx cluster.GrainContext) {
}

func (m *ManagerGrain) CreateAccount(n *proto.CreateAccountRequest, ctx cluster.GrainContext) (*proto.AccountIdResponse, error) {
	// ManagerGrainのidentityがaccountのidになる
	id := ctx.Identity()
	account := domain.NewAccount(id, n.Email)
	m.accountMap[account.ID] = account
	return &proto.AccountIdResponse{Id: account.ID}, nil
}

func (m *ManagerGrain) GetAccount(n *proto.AccountIdResponse, ctx cluster.GrainContext) (*proto.AccountResponse, error) {
	account := m.accountMap[n.Id]
	if account == nil {
		return nil, fmt.Errorf("account not found")
	}
	return &proto.AccountResponse{Id: account.ID, Email: account.Email}, nil
}

func (m *ManagerGrain) GetAllAccountEmails(n *proto.Noop, ctx cluster.GrainContext) (*proto.EmailsResponse, error) {
	emails := make(map[string]string)
	for id, account := range m.accountMap {
		emails[id] = account.Email
	}
	return &proto.EmailsResponse{Emails: emails}, nil
}