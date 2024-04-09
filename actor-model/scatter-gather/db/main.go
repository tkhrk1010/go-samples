package main

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/go-samples/go-scatter-gather-sample/db/pg"
)

func main() {
	system := actor.NewActorSystem()

	// 各FeatureCollectorActorの初期化
	tempProps := actor.PropsFromProducer(func() actor.Actor { return &TemperatureCollectorActor{} })
	tempCollector := system.Root.Spawn(tempProps)

	humidProps := actor.PropsFromProducer(func() actor.Actor { return &HumidityCollectorActor{} })
	humidCollector := system.Root.Spawn(humidProps)

	windProps := actor.PropsFromProducer(func() actor.Actor { return &WindSpeedCollectorActor{} })
	windCollector := system.Root.Spawn(windProps)

	// AggregatorActorの初期化
	aggregatorProps := actor.PropsFromProducer(func() actor.Actor {
		return &AggregatorActor{
			collectors: map[string]*actor.PID{
				"temperature": tempCollector,
				"humidity":    humidCollector,
				"windSpeed":   windCollector,
			},
		}
	})
	aggregator := system.Root.Spawn(aggregatorProps)

	// データ収集のリクエスト
	future := system.Root.RequestFuture(aggregator, &AggregateRequest{
		FeatureTypes: []string{"temperature", "humidity", "windSpeed"},
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

	db, err := pg.ConnectDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	tableName := "results"
	if err := pg.InsertData(db, tableName, response.Results); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Data inserted successfully.")

}
