package aggregator

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/go-samples/scatter-gather/fruits-trading-demo/actors/inventory/collector"
)

type AggregatorActor struct {
	Collectors map[string]*actor.PID
	results    map[string]ItemDetails
	requestor  *actor.PID
}

type AggregateRequest struct {
	ItemNames []string
}

type ItemDetails struct {
	AveragePurchasingPrice float32 `json:"averagePurchasingPrice"`
	Count                  int     `json:"count"`
}

type AggregateResponse struct {
	Results map[string]ItemDetails `json:"results"`
}

func (state *AggregatorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *AggregateRequest:
		state.results = make(map[string]ItemDetails)
		state.requestor = ctx.Sender()
		for _, itemName := range msg.ItemNames {
			clctr := state.Collectors[itemName]

			timeout := 1 * time.Second
			ctx.RequestWithCustomSender(clctr, &collector.CollectInventoryRequest{ItemName: itemName}, ctx.Self())
			ctx.SetReceiveTimeout(timeout)
		}

	case *collector.InventoryResponse:
		fmt.Println("AggregatorActor: Received InventoryResponse")
		itemDetails := ItemDetails{
			AveragePurchasingPrice: msg.AveragePurchasingPrice,
			Count:                  msg.Count,
		}
		state.results[msg.ItemName] = itemDetails

		if len(state.results) == len(state.Collectors) {
			response := &AggregateResponse{Results: state.results}
			ctx.Send(state.requestor, response)
		}
	}
}
