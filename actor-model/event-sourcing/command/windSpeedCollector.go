package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	p "github.com/tkhrk1010/protoactor-go-persistence-dynamodb/persistence"
)

//TODO: event, snapshotに、type: windSpeedCollect, windSpeedUpdatedと、id(eventId)を追加する

type WindSpeedCollectorActor struct {
	persistence.Mixin
	FeatureCollectorActor
	aggregator *actor.PID
	windSpeed float32
}

func (state *WindSpeedCollectorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CollectFeatureRequest:
		log.Println("WindSpeedCollectorActor: Received CollectFeatureRequest")
		// 最初のメッセージの送信者を保持
		state.aggregator = ctx.Sender()

		// 風速収集のロジック
		state.windSpeed = rand.Float32() * 20 // 0から20m/sの範囲でランダム

		if !state.Recovering() {
			state.PersistReceive(&p.Event{Data: fmt.Sprintf("%f", state.windSpeed)})
		}
		// Respondしても、snapshot保存時はctxが変わってるのでDeadLetterになる。
		// なので、Sendで送信する。
		ctx.Send(state.aggregator, &FeatureResponse{FeatureType: "windSpeed", Value: state.windSpeed})

	case *p.Event:
		log.Printf("Event message: %v", msg)
	case *p.Snapshot:
		log.Printf("Snapshot message: %v", msg)

	case *persistence.RequestSnapshot:
		log.Printf("RequestSnapshot message: %v", msg)
		state.PersistSnapshot(&p.Snapshot{Data: fmt.Sprintf("%f", state.windSpeed)})

	case *persistence.ReplayComplete:
		log.Printf("ReplayComplete message: %v", msg)

	case *actor.ReceiveTimeout:
		log.Println("WindSpeedCollectorActor: Received timeout")
		ctx.Send(state.aggregator, &FeatureResponse{FeatureType: "windSpeed", Value: -1.0})
	default:
		log.Printf("Unknown message: %v, message type: %T", msg, msg)
	}
}
