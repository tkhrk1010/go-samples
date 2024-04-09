package handler

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/go-samples/actor-model/scatter-gather/fruits-trading-demo/actors/market"
)

func (handler *TradeSupportInformationHandler) collectMarketInfo() (*market.Response, error) {
	marketProps := actor.PropsFromProducer(func() actor.Actor { return &market.Actor{} })

	result, err := handler.requestActorResponse(marketProps, &market.Request{
		ItemNames: []string{"apple", "orange", "banana"},
	}, 10*time.Second)
	if err != nil {
		return nil, err
	}

	mktResp, ok := result.(*market.Response)
	if !ok {
		return nil, fmt.Errorf("invalid response type")
	}
	return mktResp, nil
}
