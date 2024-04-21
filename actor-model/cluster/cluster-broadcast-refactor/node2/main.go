package main

import (
	"fmt"
	"time"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/proto"

	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"

	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/remote"
)

func main() {
	c := startNode(8081)

	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()

	fmt.Print("\nAdding 1 Egg - Enter\n")
	console.ReadLine()
	calcAdd(c, "Eggs", 1)

	fmt.Print("\nAdding 10 Egg - Enter\n")
	console.ReadLine()
	calcAdd(c, "Eggs", 10)

	fmt.Print("\nAdding 100 Bananas - Enter\n")
	console.ReadLine()
	calcAdd(c, "Bananas", 100)

	fmt.Print("\nAdding 2 Meat - Enter\n")
	console.ReadLine()
	calcAdd(c, "Meat", 3)
	calcAdd(c, "Meat", 9000)

	getAll(c)

	console.ReadLine()

	c.Shutdown(true)
}

func startNode(port int64) *cluster.Cluster {
	system := actor.NewActorSystem()

	provider := automanaged.NewWithConfig(2*time.Second, 6330, "localhost:6330", "localhost:6331")
	lookup := disthash.New()
	config := remote.Configure("localhost", 0)

	calculatorKind := proto.NewCalculatorKind(func() proto.Calculator {
		return &grain.CalcGrain{}
	}, 0)

	trackerKind := proto.NewTrackerKind(func() proto.Tracker {
		return &grain.TrackGrain{}
	}, 0)

	clusterConfig := cluster.Configure("my-cluster", provider, lookup, config,
		cluster.WithKinds(calculatorKind, trackerKind))

	cluster := cluster.New(system, clusterConfig)

	cluster.StartMember()
	return cluster
}

func calcAdd(cluster *cluster.Cluster, grainId string, addNumber int64) {
	calcGrain := proto.GetCalculatorGrainClient(cluster, grainId)
	total1, err := calcGrain.Add(&proto.NumberRequest{Number: addNumber})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Grain: %v - Total: %v \n", calcGrain.Identity, total1.Number)
}

func getAll(cluster *cluster.Cluster) {
	trackerGrain := proto.GetTrackerGrainClient(cluster, "singleTrackerGrain")
	totals, err := trackerGrain.BroadcastGetCounts(&proto.Noop{})
	if err != nil {
		panic(err)
	}

	fmt.Println("--- Totals ---")
	for grainId, count := range totals.Totals {
		fmt.Printf("%v : %v\n", grainId, count)
	}
}
