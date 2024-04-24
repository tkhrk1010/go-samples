// usecaseも含む。というかusecase寄り。分解してもいいかもしれないが、usecaseがprotoactorのGrainに依存するのも変な感じがする。
// actor model自体は、どちらかというとpresentation~infra寄りで整理されることが多いらしいし、直感的にもpresentation層(framework)寄り
// このへんのlayeringにはあまりこだわらないことにする。
package presentation

import (
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/google/uuid"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func RegisterAccount(cluster *cluster.Cluster, email string) string {
	managerGrainID := generateUUID()
	managerGrain := proto.GetManagerGrainClient(cluster, managerGrainID)
	account, err := managerGrain.CreateAccount(&proto.CreateAccountRequest{Email: email})
	if err != nil {
		panic(err)
	}

	return account.Id
}

func GetAllAccounts(cluster *cluster.Cluster, grainIds []string) map[string]string {
	var accountMap = make(map[string]string)
	for _, grainId := range grainIds {
		managerGrain := proto.GetManagerGrainClient(cluster, grainId)
		emails, err := managerGrain.GetAllAccountEmails(&proto.Noop{})
		if err != nil {
			panic(err)
		}

		for accountId, email := range emails.Emails {
			accountMap[accountId] = email
		}
	}
	return accountMap
}

func generateUUID() string {
	return uuid.New().String()
}
