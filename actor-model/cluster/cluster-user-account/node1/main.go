package main

import (
	"fmt"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/presentation"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"

	console "github.com/asynkron/goconsole"
)

// nodeをclusterに登録し、APIのようなもの(CLI)を提供する
func main() {
	proto.ManagerFactory(func() proto.Manager {
		return &grain.ManagerGrain{}
	})
	proto.AccountFactory(func() proto.Account {
		return &grain.AccountGrain{}
	})

	// node1は6330で起動。node2は6331で起動
	c := cluster.StartNode(6330)

	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()

	fmt.Print("\nAdding 1 Egg - Enter\n")
	console.ReadLine()
	// Eggsはgrain id
	presentation.CalcAdd(c, "Eggs", 1)

	fmt.Print("\nAdding 10 Egg - Enter\n")
	console.ReadLine()
	presentation.CalcAdd(c, "Eggs", 10)

	fmt.Print("\nAdding 100 Bananas - Enter\n")
	console.ReadLine()
	presentation.CalcAdd(c, "Bananas", 100)

	fmt.Print("\nAdding 2 Meat - Enter\n")
	console.ReadLine()
	presentation.CalcAdd(c, "Meat", 3)
	presentation.CalcAdd(c, "Meat", 9000)

	presentation.GetAll(c)

	console.ReadLine()

	c.Shutdown(true)
}
