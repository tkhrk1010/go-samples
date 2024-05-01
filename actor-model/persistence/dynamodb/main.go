package main

import (
	"log"
	"context"
	"fmt"
	"time"


	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	console "github.com/asynkron/goconsole"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	a "github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/actor"
	p "github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/persistence"
)

func main() {
	log.Printf("start")
	system := actor.NewActorSystem()

	client := InitializeDynamoDBClient()

	provider := p.NewProviderState(client)

	props := actor.PropsFromProducer(a.NewUserAccount, actor.WithReceiverMiddleware(persistence.Using(provider)))

	myActorPid := system.Root.Spawn(props)
	log.Printf("MyActor PID: %s", myActorPid)

	system.Root.Send(myActorPid, &p.Event{Data: "first message: please sum =+ 1"})
	time.Sleep(1 * time.Second)

	system.Root.Send(myActorPid, &persistence.RequestSnapshot{})

	_, _ = console.ReadLine()
	log.Print("done")
}

func InitializeDynamoDBClient() *dynamodb.Client {
	ctx := context.TODO()
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:4566",
				SigningRegion: "us-east-1",
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	return dynamodb.NewFromConfig(cfg)
}
