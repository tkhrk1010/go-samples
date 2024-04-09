package collector

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/tkhrk1010/go-samples/actor-model/scatter-gather/fruits-trading-demo/domain/inventory"
)

type BananaActor struct {
	InventoryCollectorActor
}

func (state *BananaActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *CollectInventoryRequest:
		fmt.Println("BananaActor: Received CollectInventoryRequest")

		// collect banana inventory info
		bananas := inventory.Banana{}
		averagePurchasingPrice := bananas.GetAveragePurchasingPrice()
		count := bananas.GetCount()

		ctx.Respond(&InventoryResponse{ItemName: "banana", AveragePurchasingPrice: averagePurchasingPrice, Count: count})

	case *actor.ReceiveTimeout:
		fmt.Println("BananaActor: Received timeout")
		ctx.Respond(&InventoryResponse{ItemName: "banana", AveragePurchasingPrice: 0.0, Count: 0})
	}
}
