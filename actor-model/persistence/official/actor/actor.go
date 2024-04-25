package actor

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/tkhrk1010/go-samples/actor-model/persistence/official/proto"
)

type Actor struct {
	persistence.Mixin
	state string
}

func (a *Actor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Println("actor started")
	case *persistence.RequestSnapshot:
		log.Printf("snapshot internal state '%v'", a.state)
		a.PersistSnapshot(&proto.Snapshot{ProtoMsg: &proto.ProtoMsg{State: a.state}})
	case *proto.Snapshot:
		a.state = msg.ProtoMsg.State
		log.Printf("recovered from snapshot, internal state changed to '%v'", a.state)
	case *persistence.ReplayComplete:
		log.Printf("replay completed, internal state changed to '%v'", a.state)
	case *proto.Message:
		scenario := "received replayed event"
		if !a.Recovering() {
			a.PersistReceive(msg)
			scenario = "received new message"
		}
		a.state = msg.ProtoMsg.State
		log.Printf("%s, internal state changed to '%v'\n", scenario, a.state)
	}
}
