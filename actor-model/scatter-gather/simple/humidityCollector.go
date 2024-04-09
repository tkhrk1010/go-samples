package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type HumidityCollectorActor struct {
	FeatureCollectorActor
}

func (state *HumidityCollectorActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *CollectFeatureRequest:
		fmt.Println("HumidityCollectorActor: Received CollectFeatureRequest")
		// 湿度収集のロジック
		rand.Seed(time.Now().UnixNano())
		humidity := rand.Float32() * 100 // 0%から100%の範囲でランダム

		ctx.Respond(&FeatureResponse{FeatureType: "humidity", Value: humidity})

	case *actor.ReceiveTimeout:
		fmt.Println("HumidityCollectorActor: Received timeout")
		ctx.Respond(&FeatureResponse{FeatureType: "humidity", Value: 1.0})
	}
}
