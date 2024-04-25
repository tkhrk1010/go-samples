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
	AmccountGrainID := generateUUID()
	AmccountGrain := proto.GetAccountGrainClient(cluster, AmccountGrainID)
	account, err := AmccountGrain.CreateAccount(&proto.CreateAccountRequest{Email: email})
	if err != nil {
		panic(err)
	}

	return account.Id
}

func GetAllAccounts(cluster *cluster.Cluster, grainIds []string) map[string]string {
	var accountMap = make(map[string]string)
	for _, grainId := range grainIds {
		AmccountGrain := proto.GetAccountGrainClient(cluster, grainId)
		res, err := AmccountGrain.GetAccountEmail(&proto.Noop{})
		if err != nil {
			panic(err)
		}
		accountMap[grainId] = res.Email
	}
	return accountMap
}

func generateUUID() string {
	return uuid.New().String()
}
