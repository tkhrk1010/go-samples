package main

import (
	"fmt"

	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	a "github.com/tkhrk1010/go-samples/actor-model/persistence/user-account/actor"
	p "github.com/tkhrk1010/go-samples/actor-model/persistence/user-account/persistence"
	"github.com/tkhrk1010/go-samples/actor-model/persistence/user-account/proto"
)

func main() {
	system := actor.NewActorSystem()
	provider := p.NewProvider(3)

	rootContext := system.Root
	props := actor.PropsFromProducer(func() actor.Actor { return &a.AccountActor{} },
		actor.WithReceiverMiddleware(persistence.Using(provider)))
	pid, _ := rootContext.SpawnNamed(props, "persistent")
	rootContext.Send(pid, &proto.AccountEvent{Msg: &proto.Event{Payload: "email1"}})
	rootContext.Send(pid, &proto.AccountEvent{Msg: &proto.Event{Payload: "email2"}})
	rootContext.Send(pid, &proto.AccountEvent{Msg: &proto.Event{Payload: "email3"}})
	rootContext.Send(pid, &proto.AccountEvent{Msg: &proto.Event{Payload: "email4"}})
	rootContext.Send(pid, &proto.AccountEvent{Msg: &proto.Event{Payload: "email5"}})

	rootContext.PoisonFuture(pid).Wait()
	fmt.Printf("*** restart ***\n")
	_, _ = rootContext.SpawnNamed(props, "persistent")

	_, _ = console.ReadLine()
}
