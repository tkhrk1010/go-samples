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
	c := cluster.StartNode("my-cluster", 6331)

	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()

	fmt.Print("\nAdding Account - Enter\n")
	console.ReadLine()
	presentation.RegisterAccount(c, "email5")

	fmt.Print("\nAdding Account - Enter\n")
	console.ReadLine()
	presentation.RegisterAccount(c, "email6")

	fmt.Print("\nAdding Account - Enter\n")
	console.ReadLine()
	presentation.RegisterAccount(c, "email7")
	presentation.RegisterAccount(c, "email8")

	console.ReadLine()

	presentation.GetAllAccounts(c)

	console.ReadLine()

	c.Shutdown(true)
}
