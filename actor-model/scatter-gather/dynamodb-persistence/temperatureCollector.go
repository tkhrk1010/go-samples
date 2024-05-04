package main

import (
	"fmt"
	"math/rand"
	"time"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	p "github.com/tkhrk1010/protoactor-go-persistence-dynamodb/persistence"
)

type TemperatureCollectorActor struct {
	persistence.Mixin
	FeatureCollectorActor
	temperature float32
}

func (state *TemperatureCollectorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CollectFeatureRequest:
		log.Println("TemperatureCollectorActor: Received CollectFeatureRequest")
		// 気温収集のロジック
		rand.Seed(time.Now().UnixNano())
		state.temperature = rand.Float32()*30 + 10 // 10°Cから40°Cの範囲でランダム

		state.PersistReceive(&p.Event{Data: fmt.Sprintf("%f", state.temperature)})
		ctx.Respond(&FeatureResponse{FeatureType: "temperature", Value: state.temperature})
	
	case *p.Snapshot:
		log.Printf("Snapshot message: %v", msg)
	case *persistence.RequestSnapshot:
		log.Printf("RequestSnapshot message: %v", msg)
		state.PersistSnapshot(&p.Snapshot{Data: fmt.Sprintf("%f", state.temperature)})

	case *persistence.ReplayComplete:
		log.Printf("ReplayComplete message: %v", msg)

	case *actor.ReceiveTimeout:
		log.Println("TemperatureCollectorActor: Received timeout")
		ctx.Respond(&FeatureResponse{FeatureType: "temperature", Value: -1.0})
	}
}
