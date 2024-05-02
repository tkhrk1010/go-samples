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
	// 基本
	system := actor.NewActorSystem()
	client := InitializeDynamoDBClient()
	provider := p.NewProviderState(client)
	props := actor.PropsFromProducer(a.NewUserAccount, actor.WithReceiverMiddleware(persistence.Using(provider)))

	userAccountActorPid := system.Root.Spawn(props)
	log.Printf("userAccountActor PID: %s", userAccountActorPid)

	system.Root.Send(userAccountActorPid, &p.Event{Data: "event1"})
	time.Sleep(1 * time.Second)

	system.Root.Send(userAccountActorPid, &persistence.RequestSnapshot{})

	_, _ = console.ReadLine()
	log.Print("done")
}

func InitializeDynamoDBClient() *dynamodb.Client {
	// Actor type *actor.UserAccount does not implement Name methodは、contextがあやしい？
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
