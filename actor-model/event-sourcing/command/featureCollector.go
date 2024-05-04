package main

import (
	"github.com/asynkron/protoactor-go/actor"
)

type FeatureCollectorActor struct{}

type CollectFeatureRequest struct {
	FeatureType string
}

type FeatureResponse struct {
	FeatureType string
	Value       float32
}

func (state *FeatureCollectorActor) Receive(ctx actor.Context) {
}
