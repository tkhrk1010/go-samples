package collector

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/go-samples/actor-model/scatter-gather/fruits-trading-demo/domain/inventory"
)

type OrangeActor struct {
	InventoryCollectorActor
}

func (state *OrangeActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *CollectInventoryRequest:
		fmt.Println("OrangeActor: Received CollectInventoryRequest")

		// collect orange inventory info
		oranges := inventory.Orange{}
		averagePurchasingPrice := oranges.GetAveragePurchasingPrice()
		count := oranges.GetCount()

		ctx.Respond(&InventoryResponse{ItemName: "orange", AveragePurchasingPrice: averagePurchasingPrice, Count: count})

	case *actor.ReceiveTimeout:
		fmt.Println("OrangeActor: Received timeout")
		ctx.Respond(&InventoryResponse{ItemName: "orange", AveragePurchasingPrice: 0.0, Count: 0})
	}
}
