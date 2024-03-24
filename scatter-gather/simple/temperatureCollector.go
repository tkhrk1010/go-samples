package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type TemperatureCollectorActor struct {
	FeatureCollectorActor
}

func (state *TemperatureCollectorActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *CollectFeatureRequest:
		fmt.Println("TemperatureCollectorActor: Received CollectFeatureRequest")
		// 気温収集のロジック
		rand.Seed(time.Now().UnixNano())
		temp := rand.Float32()*30 + 10 // 10°Cから40°Cの範囲でランダム

		ctx.Respond(&FeatureResponse{FeatureType: "temperature", Value: temp})
	case *actor.ReceiveTimeout:
		fmt.Println("TemperatureCollectorActor: Received timeout")
		ctx.Respond(&FeatureResponse{FeatureType: "temperature", Value: -1.0})
	}
}
