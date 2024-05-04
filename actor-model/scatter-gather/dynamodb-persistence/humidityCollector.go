package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	p "github.com/tkhrk1010/protoactor-go-persistence-dynamodb/persistence"
)

type HumidityCollectorActor struct {
	persistence.Mixin
	FeatureCollectorActor
	humidity float32
}

func (state *HumidityCollectorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CollectFeatureRequest:
		log.Println("HumidityCollectorActor: Received CollectFeatureRequest")
		// 湿度収集のロジック
		rand.Seed(time.Now().UnixNano())
		state.humidity = rand.Float32() * 100 // 0%から100%の範囲でランダム

		state.PersistReceive(&p.Event{Data: fmt.Sprintf("%f", state.humidity)})
		ctx.Respond(&FeatureResponse{FeatureType: "humidity", Value: state.humidity})
	case *p.Snapshot:
		log.Printf("Snapshot message: %v", msg)	
	case *persistence.RequestSnapshot:
		log.Printf("RequestSnapshot message: %v", msg)
		state.PersistSnapshot(&p.Snapshot{Data: fmt.Sprintf("%f", state.humidity)})
	
	case *persistence.ReplayComplete:
		log.Printf("ReplayComplete message: %v", msg)
	
	case *actor.ReceiveTimeout:
		log.Println("HumidityCollectorActor: Received timeout")
		ctx.Respond(&FeatureResponse{FeatureType: "humidity", Value: 1.0})
	}
}
