package handler

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/go-samples/scatter-gather/fruits-trading-demo/actors/inventory/aggregator"
	"github.com/tkhrk1010/go-samples/scatter-gather/fruits-trading-demo/actors/inventory/collector"
)

func (handler *TradeSupportInformationHandler) collectInventoryInfo() (*aggregator.AggregateResponse, error) {
	orangeInventoryCollector := handler.spawnActor(func() actor.Actor { return &collector.OrangeActor{} })
	appleInventoryCollector := handler.spawnActor(func() actor.Actor { return &collector.AppleActor{} })
	bananaInventoryCollector := handler.spawnActor(func() actor.Actor { return &collector.BananaActor{} })

	aggregatorProps := actor.PropsFromProducer(func() actor.Actor {
		return &aggregator.AggregatorActor{
			Collectors: map[string]*actor.PID{
				"appleInventory":  appleInventoryCollector,
				"orangeInventory": orangeInventoryCollector,
				"bananaInventory": bananaInventoryCollector,
			},
		}
	})

	result, err := handler.requestActorResponse(aggregatorProps, &aggregator.AggregateRequest{
		ItemNames: []string{"orangeInventory", "appleInventory", "bananaInventory"},
	}, 10*time.Second)
	if err != nil {
		return nil, err
	}

	aggResp, ok := result.(*aggregator.AggregateResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response type")
	}
	return aggResp, nil
}
