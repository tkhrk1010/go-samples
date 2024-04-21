package presentation

import (
	"fmt"

	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/proto"
)

func CalcAdd(cluster *cluster.Cluster, grainId string, addNumber int64) {
	itemGrain := proto.GetItemGrainClient(cluster, grainId)
	total1, err := itemGrain.Add(&proto.NumberRequest{Number: addNumber})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Grain: %v - Total: %v \n", itemGrain.Identity, total1.Number)
}

func GetAll(cluster *cluster.Cluster) {
	cartGrain := proto.GetCartGrainClient(cluster, "singleCartGrain")
	totals, err := cartGrain.BroadcastGetCounts(&proto.Noop{})
	if err != nil {
		panic(err)
	}

	fmt.Println("--- Totals ---")
	for grainId, count := range totals.Totals {
		fmt.Printf("%v : %v\n", grainId, count)
	}
}
