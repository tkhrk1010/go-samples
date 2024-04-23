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
	c := cluster.StartNode("my-cluster", 6330)

	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()

	fmt.Print("\nAdding Account - Enter\n")
	console.ReadLine()
	presentation.RegisterAccount(c, "email1")

	fmt.Print("\nAdding Account - Enter\n")
	console.ReadLine()
	presentation.RegisterAccount(c, "email2")

	fmt.Print("\nAdding Account - Enter\n")
	console.ReadLine()
	presentation.RegisterAccount(c, "email3")
	presentation.RegisterAccount(c, "email4")

	console.ReadLine()

	presentation.GetAllAccounts(c)

	console.ReadLine()

	c.Shutdown(true)
}
