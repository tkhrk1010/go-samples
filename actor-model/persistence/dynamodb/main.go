package main

import (
	"log"
	"context"
	"fmt"
	"time"
	"errors"


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

	userAccount1, err := system.Root.SpawnNamed(props, "userAccountActor-"+"1")
	// 登録ユーザーのメールアドレスが既に存在する場合はエラーを返す
	// メッセージ送信時に現在のバージョンを送信することで、永続化されたデータとの競合を防ぐことができるらしい
	// 詳しくはprotobufを参照してください
	// TODO: protobufを見て勉強する
	// Ref: github.com/ytake/protoactor-go-cqrs-example/internal/registration/create_user.go
	if errors.Is(err, actor.ErrNameExists) {
		log.Printf("user %s already exists", userAccount1)
	}
	if err != nil {
		log.Printf("failed error %s", err.Error())
	}
	log.Printf("userAccountActor PID: %s", userAccount1)

	system.Root.Send(userAccount1, &p.Event{Data: "event1"})
	system.Root.Send(userAccount1, &p.Event{Data: "event2"})
	system.Root.Send(userAccount1, &p.Event{Data: "event3"})
	system.Root.Send(userAccount1, &p.Event{Data: "event4"})
	time.Sleep(1 * time.Second)

	// 同じactorNameのactorが生まれたらerrorになることを確認
	// sameUserAccount1, err := system.Root.SpawnNamed(props, "userAccountActor-"+"1") 
	// if errors.Is(err, actor.ErrNameExists) {
	// 	log.Printf("user %s already exists", sameUserAccount1)
	// }
	// if err != nil {
	// 	log.Printf("failed error %s", err.Error())
	// }
	// log.Printf("userAccountActor PID: %s", sameUserAccount1)

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
