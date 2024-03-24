package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/go-samples/scatter-gather/fruits-trading-demo/actors/inventory/aggregator"
	"github.com/tkhrk1010/go-samples/scatter-gather/fruits-trading-demo/actors/market"
)

type TradeSupportInformationHandler struct {
	system *actor.ActorSystem
}

type CombinedItemDetails struct {
	AveragePurchasingPrice float32 `json:"averagePurchasingPrice"`
	Count                  int     `json:"count"`
	MarketPrice            float32 `json:"marketPrice"`
}

type CombinedResponse struct {
	Items map[string]CombinedItemDetails `json:"items"`
}

func (handler *TradeSupportInformationHandler) GetTradeSupportInformation() ([]byte, error) {
	aggregatorResponse, err := handler.collectInventoryInfo()
	if err != nil {
		return nil, err
	}

	marketResponse, err := handler.collectMarketInfo()
	if err != nil {
		return nil, err
	}

	combinedResponse := CombineResponses(aggregatorResponse, marketResponse)

	jsonResponse, err := json.Marshal(combinedResponse)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil, err
	}

	return jsonResponse, nil
}

func NewTradeSupportInformationHandler() *TradeSupportInformationHandler {
	return &TradeSupportInformationHandler{
		system: actor.NewActorSystem(),
	}
}

func (handler *TradeSupportInformationHandler) spawnActor(producer func() actor.Actor) *actor.PID {
	props := actor.PropsFromProducer(producer)
	return handler.system.Root.Spawn(props)
}

func (handler *TradeSupportInformationHandler) requestActorResponse(props *actor.Props, request interface{}, timeout time.Duration) (interface{}, error) {
	actorRef := handler.system.Root.Spawn(props)
	future := handler.system.Root.RequestFuture(actorRef, request, timeout)
	return future.Result()
}

func CombineResponses(aggregatorResponse *aggregator.AggregateResponse, marketResponse *market.Response) CombinedResponse {
	combined := CombinedResponse{
		Items: make(map[string]CombinedItemDetails),
	}

	for itemName, aggDetails := range aggregatorResponse.Results {
		combined.Items[itemName] = CombinedItemDetails{
			AveragePurchasingPrice: aggDetails.AveragePurchasingPrice,
			Count:                  aggDetails.Count,
		}
	}

	for itemName, marketItemDetails := range marketResponse.Results {
		if item, exists := combined.Items[itemName]; exists {
			item.MarketPrice = marketItemDetails.MarketPrice
			combined.Items[itemName] = item
		}
	}

	return combined
}
