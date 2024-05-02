package actor

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	p "github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/persistence"
)

// Nameというfiledを持ってしまうと、MixinのNameと競合してしまい、エラーになるので注意
type UserAccount struct {
	persistence.Mixin
	Id    string
	Email string
}

func (u *UserAccount) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Printf("actor started: %v", ctx.Self())
	case *p.Event:
		log.Printf("Event message: %v", msg)
	case *p.Snapshot:
		log.Printf("Snapshot message: %v", msg)
	case *persistence.RequestSnapshot:
		log.Printf("RequestSnapshot message: %v", msg)
	case *persistence.ReplayComplete:
		log.Printf("ReplayComplete message: %v", msg)
	default:
		log.Printf("Unknown message: %v, message type: %T", msg, msg)
	}
}

func NewUserAccount() actor.Actor {
	return &UserAccount{}
}
