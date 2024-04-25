package main

import (
	"fmt"

	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	a "github.com/tkhrk1010/go-samples/actor-model/persistence/user-account/actor"
	p "github.com/tkhrk1010/go-samples/actor-model/persistence/user-account/persistence"
)

func main() {
	system := actor.NewActorSystem()
	provider := p.NewProvider(3)

	rootContext := system.Root
	props := actor.PropsFromProducer(func() actor.Actor { return &a.AccountActor{} },
		actor.WithReceiverMiddleware(persistence.Using(provider)))
	pid1, _ := rootContext.SpawnNamed(props, "user-1")
	pid2, _ := rootContext.SpawnNamed(props, "user-2")
	pid3, _ := rootContext.SpawnNamed(props, "user-3")
	pid4, _ := rootContext.SpawnNamed(props, "user-4")
	pid5, _ := rootContext.SpawnNamed(props, "user-5")

	// RFC3339とは、ISO 8601に基づいた日付と時刻の表現方法の一つで、yyyy-mm-ddThh:mm:ss:msの形式で表現される
	// 例: 2021-01-01T00:00:00.000Z
	rootContext.Send(pid1, a.NewSignUpEvent("email1"))
	rootContext.Send(pid2, a.NewSignUpEvent("email2"))
	rootContext.Send(pid3, a.NewSignUpEvent("email3"))
	rootContext.Send(pid4, a.NewSignUpEvent("email4"))
	rootContext.Send(pid5, a.NewSignUpEvent("email5"))

	rootContext.Send(pid1, a.NewSignUpEvent("email1-2"))
	rootContext.Send(pid1, a.NewSignUpEvent("email1-3"))
	rootContext.Send(pid1, a.NewSignUpEvent("email1-4"))
	rootContext.Send(pid1, a.NewSignUpEvent("email1-5"))

	rootContext.PoisonFuture(pid1).Wait()
	fmt.Printf("*** restart ***\n")
	_, _ = rootContext.SpawnNamed(props, "user-1")

	_, _ = console.ReadLine()
}
