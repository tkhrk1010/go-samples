package grain

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
	a "github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/actor"

)

type ManagerGrain struct {
	accountActor *actor.PID
}

func (m *ManagerGrain) ReceiveDefault(ctx cluster.GrainContext) {
	log.Printf("ManagerGrain %s received: %v", ctx.Identity(), ctx.Message())
}

func (m *ManagerGrain) Init(ctx cluster.GrainContext) {
	log.Printf("ManagerGrain %s initialized", ctx.Identity())
	m.accountActor = ctx.Spawn(actor.PropsFromProducer(func() actor.Actor { return a.NewAccountActor(ctx.Identity()) }))
}

func (m *ManagerGrain) Terminate(ctx cluster.GrainContext) {
	log.Printf("ManagerGrain %s terminated", ctx.Identity())
}

func (m *ManagerGrain) CreateAccount(n *proto.CreateAccountRequest, ctx cluster.GrainContext) (*proto.AccountIdResponse, error) {
	// ManagerGrainのidentityがaccountのidになる
	id := ctx.Identity()
	accountActor := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor { return a.NewAccountActor() }))
	m.accountActor = accountActor

	future := ctx.RequestFuture(accountActor, &proto.CreateAccountRequest{Id: id, Email: n.Email}, 30*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Printf("CreateAccountRequest grainToActor error: %v", err)
		return nil, err
	}
	return result.(*proto.AccountIdResponse), nil
}

func (m *ManagerGrain) GetAccount(n *proto.AccountIdResponse, ctx cluster.GrainContext) (*proto.AccountResponse, error) {
	future := ctx.RequestFuture(m.accountActor, &proto.Noop{}, 30*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Printf("GetAccount grainToActor error: %v", err)
		return nil, err
	}
	return result.(*proto.AccountResponse), nil
}


func (m *ManagerGrain) GetAccountEmail(n *proto.Noop, ctx cluster.GrainContext) (*proto.EmailResponse, error) {
	log.Printf("ManagerGrain %s GetAccountEmail", ctx.Identity())
	future := ctx.RequestFuture(m.accountActor, &proto.Noop{}, 100*time.Millisecond)
	result, err := future.Result()
	if err != nil {
		log.Printf("Error GetAccountEmail: %v", err)
		panic(err)
	}
	account := result.(*proto.AccountResponse)
	if account == nil {
		log.Printf("Error GetAccountEmail: %v", "account is nil")
		panic("account is nil")
	}
	if account.Email == "" {
		log.Printf("Error GetAccountEmail: %v", "email is empty")
		panic("grain email is empty")
	}

	return &proto.EmailResponse{Email: account.Email}, nil
}
