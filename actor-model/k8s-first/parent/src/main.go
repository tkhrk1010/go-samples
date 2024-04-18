package main

import (
	ctx "context"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/go-samples/actor-model/k8s-first/parent/src/proto"
	"google.golang.org/grpc"
)

type ParentActor struct {
	childPID *actor.PID
}

func NewParentActor() actor.Actor {
	return &ParentActor{}
}

func (state *ParentActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Println("Parent actor started")
		log.Println("Parent actor messege typed *actor.Started: ", msg)

	case CreateChild:
		log.Println("Parent actor received CreateChild message")
		conn, err := grpc.Dial("child-service:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to child service: %v", err)
		}
		defer conn.Close()

		client := proto.NewActorServiceClient(conn)
		resp, err := client.CreateChild(ctx.Background(), &proto.CreateChildRequest{})
		if err != nil {
			log.Fatalf("Failed to create child actor: %v", err)
		}
		log.Printf("Child actor created with ID: %s", resp.ChildId)
	}
}

type CreateChild struct{}

func main() {
	log.SetFlags(log.Lmicroseconds)

	log.Println("start")
	system := actor.NewActorSystem()

	parentProps := actor.PropsFromProducer(NewParentActor)
	parentPID := system.Root.Spawn(parentProps)
	log.Printf("Parent actor PID: %v\n", parentPID)

	log.Println("send CreateChild message to ParentActor")
	system.Root.Send(parentPID, CreateChild{})

	console := make(chan struct{})
	<-console
}
