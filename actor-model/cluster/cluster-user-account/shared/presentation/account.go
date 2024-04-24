package presentation

import (
	"fmt"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func RegisterAccount(cluster *cluster.Cluster, email string) string {
	managerGrain := proto.GetManagerGrainClient(cluster, "singleManagerGrain")
	account, err := managerGrain.CreateAccount(&proto.CreateAccountRequest{Email: email})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Account registered. ID: %v \n", account.Id)
	return account.Id
}

func GetAllAccounts(cluster *cluster.Cluster) {
	managerGrain := proto.GetManagerGrainClient(cluster, "singleManagerGrain")
	emails, err := managerGrain.GetAllAccountEmails(&proto.Noop{})
	if err != nil {
		panic(err)
	}

	fmt.Println("--- Emails ---")
	for accountId, email := range emails.Emails {
		fmt.Printf("%v : %v\n", accountId, email)
	}
}
