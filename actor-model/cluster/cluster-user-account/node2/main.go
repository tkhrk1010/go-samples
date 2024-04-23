package main

import (
	"fmt"

	console "github.com/asynkron/goconsole"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/presentation"
)

// nodeをclusterに登録し、APIのようなもの(CLI)を提供する
func main() {
	// Factoryの登録はnode1で行っているため、node2では不要

	// node1は6330で起動。node2は6331で起動
	c := cluster.StartNode(6331)

	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()

	fmt.Print("\nAdding 1 Egg - Enter\n")
	console.ReadLine()
	presentation.RegisterAccount(c, "Eggs", 1)

	fmt.Print("\nAdding 10 Egg - Enter\n")
	console.ReadLine()
	presentation.RegisterAccount(c, "Eggs", 10)

	fmt.Print("\nAdding 100 Bananas - Enter\n")
	console.ReadLine()
	presentation.RegisterAccount(c, "Bananas", 100)

	fmt.Print("\nAdding 2 Meat - Enter\n")
	console.ReadLine()
	presentation.RegisterAccount(c, "Meat", 3)
	presentation.RegisterAccount(c, "Meat", 9000)

	presentation.GetAllAccounts(c)

	console.ReadLine()

	c.Shutdown(true)
}
