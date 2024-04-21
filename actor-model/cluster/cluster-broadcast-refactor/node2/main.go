package main

import (
	"fmt"
	console "github.com/asynkron/goconsole"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/presentation"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-broadcast-refactor/shared/cluster"
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

