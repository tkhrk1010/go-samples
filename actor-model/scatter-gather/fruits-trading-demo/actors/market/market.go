package market

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	dm "github.com/tkhrk1010/go-samples/actor-model/scatter-gather/fruits-trading-demo/domain/market"
)

type Actor struct {
	results   map[string]ItemDetails
	requestor *actor.PID
}

type Request struct {
	ItemNames []string
}

type ItemDetails struct {
	MarketPrice float32 `json:"marketPrice"`
}

type Response struct {
	Results map[string]ItemDetails `json:"results"`
}

func (state *Actor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Request:
		state.results = make(map[string]ItemDetails)
		state.requestor = ctx.Sender()
		fmt.Println("market Actor: Received Request")

		// api call to market (dummy)
		apiResponse := map[string]float32{
			"apple":  100.0,
			"orange": 200.0,
			"banana": 300.0,
			"grape":  400.0,
			"melon":  500.0,
		}

		for _, itemName := range msg.ItemNames {
			var model dm.MarketItem
			switch itemName {
			case "apple":
				model = &dm.Apple{
					Price: apiResponse[itemName],
				}
			case "orange":
				model = &dm.Orange{
					Price: apiResponse[itemName],
				}
			case "banana":
				model = &dm.Banana{
					Price: apiResponse[itemName],
				}
			default:
				fmt.Println("Actor: invalid itemName")
				continue
			}

			itemDetails := ItemDetails{
				MarketPrice: model.GetPrice(),
			}
			state.results[itemName] = itemDetails
		}

		response := &Response{Results: state.results}
		ctx.Send(state.requestor, response)
	}
}
