package main

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"

	p "github.com/tkhrk1010/protoactor-go-persistence-dynamodb/persistence"
)

func main() {
	system := actor.NewActorSystem()
	client := p.InitializeDynamoDBClient()
	provider := p.NewProviderState(client)

	// 各FeatureCollectorActorの初期化
	windProps := actor.PropsFromProducer(func() actor.Actor { return &WindSpeedCollectorActor{} }, actor.WithReceiverMiddleware(persistence.Using(provider)))
	windCollector, err := system.Root.SpawnNamed(windProps, "windSpeed")
	if err != nil {
		fmt.Printf("Error while spawning actor: %s\n", err)
	}

	// AggregatorActorの初期化
	aggregatorProps := actor.PropsFromProducer(func() actor.Actor {
		return &AggregatorActor{
			collectors: map[string]*actor.PID{
				"windSpeed":   windCollector,
			},
		}
	})
	aggregator := system.Root.Spawn(aggregatorProps)

	// データ収集のリクエスト
	future := system.Root.RequestFuture(aggregator, &AggregateRequest{
		FeatureTypes: []string{"windSpeed"},
	}, 10*time.Second)

	// レスポンスの待機
	result, err := future.Result()
	if err != nil {
		fmt.Printf("Error while waiting for result: %s\n", err)
		return
	}

	// 結果の出力
	response, ok := result.(*AggregateResponse)
	if !ok {
		fmt.Println("Invalid response type")
		return
	}

	fmt.Println("Aggregate Results:", response.Results)

}
