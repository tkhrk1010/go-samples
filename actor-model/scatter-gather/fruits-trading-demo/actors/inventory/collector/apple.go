package collector

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/go-samples/scatter-gather/fruits-trading-demo/domain/inventory"
)

type AppleActor struct {
	InventoryCollectorActor
}

func (state *AppleActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *CollectInventoryRequest:
		fmt.Println("AppleActor: Received CollectInventoryRequest")

		// collect apple inventory info
		apples := inventory.Apple{}
		averagePurchasingPrice := apples.GetAveragePurchasingPrice()
		count := apples.GetCount()

		ctx.Respond(&InventoryResponse{ItemName: "apple", AveragePurchasingPrice: averagePurchasingPrice, Count: count})

	case *actor.ReceiveTimeout:
		fmt.Println("AppleActor: Received timeout")
		ctx.Respond(&InventoryResponse{ItemName: "apple", AveragePurchasingPrice: 0.0, Count: 0})
	}
}
