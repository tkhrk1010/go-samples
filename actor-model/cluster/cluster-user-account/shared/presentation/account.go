package presentation

import (
	"fmt"

	"github.com/asynkron/protoactor-go/cluster"
	"github.com/google/uuid"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func generateUUID() string {
	return uuid.New().String()
}

func RegisterAccount(cluster *cluster.Cluster, email string) string {
	grainId := generateUUID()
	accountGrain := proto.GetAccountGrainClient(cluster, grainId)
	account, err := accountGrain.Register(&proto.AccountRegisterRequest{Id: grainId, Email: email})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Grain: %v - AccountID: %v \n", accountGrain.Identity, account.Id)
	return grainId
}

func GetAllAccounts(cluster *cluster.Cluster) {
	managerGrain := proto.GetManagerGrainClient(cluster, "singleManagerGrain")
	emails, err := managerGrain.GetAllAccountEmails(&proto.Noop{})
	if err != nil {
		panic(err)
	}

	fmt.Println("--- Emails ---")
	for grainId, email := range emails.Emails {
		fmt.Printf("%v : %v\n", grainId, email)
	}
}
