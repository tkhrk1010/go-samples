package main

import (
	"fmt"

	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	a "github.com/tkhrk1010/go-samples/actor-model/persistence/official/actor"
	"github.com/tkhrk1010/go-samples/actor-model/persistence/official/proto"
	p "github.com/tkhrk1010/go-samples/actor-model/persistence/official/persistence"
)

func main() {
	system := actor.NewActorSystem()
	provider := p.NewProvider(3)

	rootContext := system.Root
	props := actor.PropsFromProducer(func() actor.Actor { return &a.Actor{} },
		actor.WithReceiverMiddleware(persistence.Using(provider)))
	pid, _ := rootContext.SpawnNamed(props, "persistent")
	rootContext.Send(pid, &proto.Message{ProtoMsg: &proto.ProtoMsg{State: "state1"}})
	rootContext.Send(pid, &proto.Message{ProtoMsg: &proto.ProtoMsg{State: "state2"}})
	rootContext.Send(pid, &proto.Message{ProtoMsg: &proto.ProtoMsg{State: "state3"}})
	rootContext.Send(pid, &proto.Message{ProtoMsg: &proto.ProtoMsg{State: "state4"}})
	rootContext.Send(pid, &proto.Message{ProtoMsg: &proto.ProtoMsg{State: "state5"}})

	rootContext.PoisonFuture(pid).Wait()
	fmt.Printf("*** restart ***\n")
	_, _ = rootContext.SpawnNamed(props, "persistent")

	_, _ = console.ReadLine()
}
