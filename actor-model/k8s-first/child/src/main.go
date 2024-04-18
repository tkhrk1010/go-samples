package main

import (
	"context"
	"log"
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"google.golang.org/grpc"

	"github.com/tkhrk1010/go-samples/actor-model/k8s-first/child/src/proto"
)

type ChildActor struct{}

type actorServer struct {
	proto.UnimplementedActorServiceServer
}

func (s *actorServer) CreateChild(ctx context.Context, req *proto.CreateChildRequest) (*proto.CreateChildResponse, error) {
	childProps := actor.PropsFromProducer(func() actor.Actor { return &ChildActor{} })
	childPID := actor.NewActorSystem().Root.Spawn(childProps)
	return &proto.CreateChildResponse{ChildId: childPID.Id}, nil
}

func (state *ChildActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Println("Child actor started")
		log.Println("Child actor messege typed *actor.Started: ", msg)
		log.Println("Child actor's Address: ", context.Self().Address, ", PID: ", context.Self().Id)
	}
}

func main() {
	log.SetFlags(log.Lmicroseconds)

	log.Println("start")
	actor.NewActorSystem()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterActorServiceServer(s, &actorServer{})
	log.Println("gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	console := make(chan struct{})
	<-console
}
