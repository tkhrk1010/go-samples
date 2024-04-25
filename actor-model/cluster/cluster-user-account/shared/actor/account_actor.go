package actor

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/domain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

type AccountActor struct {
	account *domain.Account
	id 		string
}

func NewAccountActor(id string) *AccountActor {
	return &AccountActor{id: id}
}

func (a *AccountActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *proto.CreateAccountRequest:
		a.account = domain.NewAccount(msg.Id, msg.Email)
		ctx.Respond(&proto.AccountIdResponse{Id: a.account.ID})
	case *proto.AccountIdResponse:
		a.account.ID = msg.Id
	case *proto.AccountResponse:
		a.account.ID = msg.Id
		a.account.Email = msg.Email
	// TODO: messageの型変える
	case *proto.Noop:
		log.Printf("AccountActor %s received: %v", ctx.Self().Id, msg)
		ctx.Respond(&proto.AccountResponse{Id: a.account.ID, Email: a.account.Email})
	}
}

