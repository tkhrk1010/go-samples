package main

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type AggregatorActor struct {
	collectors map[string]*actor.PID
	results    map[string]float32
	requestor  *actor.PID
}

type AggregateRequest struct {
	FeatureTypes []string
}

type AggregateResponse struct {
	Results map[string]float32
}


func (state *AggregatorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *AggregateRequest:
		state.results = make(map[string]float32)
		state.requestor = ctx.Sender()
		for _, featureType := range msg.FeatureTypes {
			collector := state.collectors[featureType]

			timeout := 10 * time.Second
			ctx.RequestWithCustomSender(collector, &CollectFeatureRequest{FeatureType: featureType}, ctx.Self())
			ctx.SetReceiveTimeout(timeout)
		}
	case *FeatureResponse:
		fmt.Println("AggregatorActor: Received FeatureResponse")
		state.results[msg.FeatureType] = msg.Value
		if len(state.results) == len(state.collectors) {
			response := &AggregateResponse{Results: state.results}
			ctx.Send(state.requestor, response)
			// これだと動かないので注意
			// ctx.Respond(response)
		}
	}
}
