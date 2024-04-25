package grain

import (
	"fmt"

	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/domain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

type AccountGrain struct {
	account *domain.Account
}

func (a *AccountGrain) ReceiveDefault(ctx cluster.GrainContext) {
}

func (a *AccountGrain) Init(ctx cluster.GrainContext) {
	a.account = &domain.Account{}
}

func (a *AccountGrain) Terminate(ctx cluster.GrainContext) {
}

func (a *AccountGrain) CreateAccount(n *proto.CreateAccountRequest, ctx cluster.GrainContext) (*proto.AccountIdResponse, error) {
	// AccountGrainのidentityがaccountのidになる
	id := ctx.Identity()
	a.account = domain.NewAccount(id, n.Email)
	return &proto.AccountIdResponse{Id: a.account.ID}, nil
}

func (a *AccountGrain) GetAccount(n *proto.AccountIdResponse, ctx cluster.GrainContext) (*proto.AccountResponse, error) {
	account := a.account
	if account == nil {
		return nil, fmt.Errorf("account not found")
	}
	return &proto.AccountResponse{Id: account.ID, Email: account.Email}, nil
}

func (a *AccountGrain) GetAccountEmail(n *proto.Noop, ctx cluster.GrainContext) (*proto.AccountEmailResponse, error) {
	return &proto.AccountEmailResponse{Email: a.account.Email}, nil
}
