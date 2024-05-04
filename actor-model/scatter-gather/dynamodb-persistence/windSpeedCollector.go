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

type WindSpeedCollectorActor struct {
	persistence.Mixin
	FeatureCollectorActor
	windSpeed float32
}

func (state *WindSpeedCollectorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CollectFeatureRequest:
		log.Println("WindSpeedCollectorActor: Received CollectFeatureRequest")
		// 風速収集のロジック
		rand.Seed(time.Now().UnixNano())
		state.windSpeed = rand.Float32() * 20 // 0から20m/sの範囲でランダム

		state.PersistReceive(&p.Event{Data: fmt.Sprintf("%f", state.windSpeed)})
		ctx.Respond(&FeatureResponse{FeatureType: "windSpeed", Value: state.windSpeed})
	case *p.Snapshot:
		log.Printf("Snapshot message: %v", msg)
	case *persistence.RequestSnapshot:
		log.Printf("RequestSnapshot message: %v", msg)
		state.PersistSnapshot(&p.Snapshot{Data: fmt.Sprintf("%f", state.windSpeed)})
	
	case *persistence.ReplayComplete:
		log.Printf("ReplayComplete message: %v", msg)
	
	case *actor.ReceiveTimeout:
		log.Println("WindSpeedCollectorActor: Received timeout")
		ctx.Respond(&FeatureResponse{FeatureType: "windSpeed", Value: -1.0})
	}
}
