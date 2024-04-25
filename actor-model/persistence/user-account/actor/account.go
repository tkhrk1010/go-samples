package actor

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/tkhrk1010/go-samples/actor-model/persistence/user-account/proto"
)

type AccountActor struct {
	persistence.Mixin
	// id 		string
	email 	string
}

func (a *AccountActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Println("actor started")
	case *persistence.RequestSnapshot:
		log.Printf("snapshot actor email '%v'", a.email)
		a.PersistSnapshot(&proto.AccountSnapshot{Msg: &proto.Event{Payload: a.email}})
	case *proto.AccountSnapshot:
		a.email = msg.Msg.Payload
		log.Printf("recovered from snapshot, actor email changed to '%v'", a.email)
	case *persistence.ReplayComplete:
		log.Printf("replay completed, actor email changed to '%v'", a.email)
	case *proto.AccountEvent:
		scenario := "received replayed event"
		if !a.Recovering() {
			a.PersistReceive(msg)
			scenario = "received new message"
		}
		a.email = msg.Msg.Payload
		log.Printf("%s, actor email changed to '%v'\n", scenario, a.email)
	}
}
