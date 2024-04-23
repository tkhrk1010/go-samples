package presentation

import (
	"fmt"

	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func CalcAdd(cluster *cluster.Cluster, grainId string, addNumber int64) {
	accountGrain := proto.GetAccountGrainClient(cluster, grainId)
	total1, err := accountGrain.Add(&proto.NumberRequest{Number: addNumber})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Grain: %v - Total: %v \n", accountGrain.Identity, total1.Number)
}

func GetAll(cluster *cluster.Cluster) {
	managerGrain := proto.GetManagerGrainClient(cluster, "singleManagerGrain")
	totals, err := managerGrain.BroadcastGetCounts(&proto.Noop{})
	if err != nil {
		panic(err)
	}

	fmt.Println("--- Totals ---")
	for grainId, count := range totals.Totals {
		fmt.Printf("%v : %v\n", grainId, count)
	}
}
