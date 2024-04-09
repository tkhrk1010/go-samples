package collector

import (
	"github.com/asynkron/protoactor-go/actor"
)

type InventoryCollectorActor struct{}

type CollectInventoryRequest struct {
	ItemName string
}

type InventoryResponse struct {
	ItemName               string
	AveragePurchasingPrice float32
	Count 								 int
}

func (state *InventoryCollectorActor) Receive(ctx actor.Context) {
}
