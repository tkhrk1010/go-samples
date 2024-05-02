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
		// Persist all events received outside of recovery
		if !u.Recovering() {
			u.PersistReceive(msg)
		}
		// Set state to whatever message says
		u.Email = msg.Data
	case *persistence.RequestSnapshot:
		log.Printf("RequestSnapshot message: %v", msg)
		u.PersistSnapshot(newSnapshot(u.Email))
	case *persistence.ReplayComplete:
		log.Printf("ReplayComplete message: %v", msg)
	default:
		log.Printf("Unknown message: %v, message type: %T", msg, msg)
	}
}

func NewUserAccount() actor.Actor {
	return &UserAccount{}
}

func newSnapshot(data string) *p.Snapshot {
	return &p.Snapshot{Data: data}
}
