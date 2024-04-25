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
	proto.AccountFactory(func() proto.Account {
		return &grain.AccountGrain{}
	})
	// node1は6330で起動。node2は6331で起動
	c := cluster.StartNode("my-cluster", 6330)

	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()

	fmt.Print("\nAdding Account - Enter\n")
	console.ReadLine()
	id1 := presentation.RegisterAccount(c, "email1")
	fmt.Printf("Account registered. ID: %v \n", id1)

	fmt.Print("\nAdding Account - Enter\n")
	console.ReadLine()
	id2 := presentation.RegisterAccount(c, "email2")
	fmt.Printf("Account registered. ID: %v \n", id2)

	fmt.Print("\nAdding Account - Enter\n")
	console.ReadLine()
	id3 := presentation.RegisterAccount(c, "email3")
	fmt.Printf("Account registered. ID: %v \n", id3)
	id4 := presentation.RegisterAccount(c, "email4")
	fmt.Printf("Account registered. ID: %v \n", id4)

	console.ReadLine()

	accounts := presentation.GetAllAccounts(c, []string{id1, id2, id3, id4})
	fmt.Println("--- Emails ---")
	for accountId, email := range accounts {
		fmt.Printf("%v : %v\n", accountId, email)
	}

	console.ReadLine()

	c.Shutdown(true)
}
