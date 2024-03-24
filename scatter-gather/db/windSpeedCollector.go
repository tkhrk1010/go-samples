package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type WindSpeedCollectorActor struct {
	FeatureCollectorActor
}

func (state *WindSpeedCollectorActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *CollectFeatureRequest:
		fmt.Println("WindSpeedCollectorActor: Received CollectFeatureRequest")
		// 風速収集のロジック
		rand.Seed(time.Now().UnixNano())
		windSpeed := rand.Float32() * 20 // 0から20m/sの範囲でランダム

		ctx.Respond(&FeatureResponse{FeatureType: "windSpeed", Value: windSpeed})
	case *actor.ReceiveTimeout:
		fmt.Println("WindSpeedCollectorActor: Received timeout")
		ctx.Respond(&FeatureResponse{FeatureType: "windSpeed", Value: -1.0})
	}
}
