// usecaseも含む。というかusecase寄り。分解してもいいかもしれないが、usecaseがprotoactorのGrainに依存するのも変な感じがする。
// actor model自体は、どちらかというとpresentation~infra寄りで整理されることが多いらしいし、直感的にもpresentation層(framework)寄り
// このへんのlayeringにはあまりこだわらないことにする。
package presentation

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/cluster"
	"github.com/google/uuid"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func RegisterAccount(cluster *cluster.Cluster, email string) (string, error) {
	accountId := generateUUID()
	// log.Printf("accountId: %v is registering", accountId)
	managerGrain := proto.GetManagerGrainClient(cluster, accountId)
	account, err := managerGrain.CreateAccount(&proto.CreateAccountRequest{Id: accountId, Email: email})
	if err != nil {
		log.Printf("pre: error: %v", err)
		panic(err)
	}
	// log.Printf("accountId: %v succeeded to register", account.Id)
	// res, err := managerGrain.GetAccountEmail(&proto.Noop{})
	// if err != nil {
	// 	log.Printf("pre: error in post: %v", err)
	// 	panic(err)
	// }
	// log.Printf("pre: registered accountId: %v email: %v", account.Id, res.Email)
	return account.Id, err
}

func GetAllAccounts(nodes []*cluster.Cluster, grainIds []string) map[string]string {
	cluster := nodes[0]
	var accountMap = make(map[string]string)
	for _, grainId := range grainIds {
		log.Printf("pre: try to get email from grainId: %v", grainId)
		var EmailResponse *proto.EmailResponse
		var err error
		managerGrain := proto.GetManagerGrainClient(cluster, grainId)

		// grainの取得を少し待ってみる gossip intervalが300msなので、それより長めに待つ
		// cluster/config.go
		time.Sleep(time.Millisecond * 400)

		// リトライ回数を定義
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			EmailResponse, err = managerGrain.GetAccountEmail(&proto.Noop{})
			if err == nil {
				log.Printf("pre: GetAllAccounts from grain succeeded (attempt %d), grainId: %v", i+1, grainId)
				break
			}
			log.Printf("pre: Error GetAllAccounts from grain (attempt %d): %v", i+1, err)
			// 少し待ってからリトライ
			time.Sleep(time.Second * 10)
			// 次のnodeで上書きしてみる
			cluster = nodes[1]
		}

		if err != nil {
			log.Printf("Error GetAllAccounts from grain after %d attempts: %v", maxRetries, err)
			panic(err)
		}

		if EmailResponse == nil {
			log.Printf("error GetAllAccounts EmailResponse nil")
			panic("EmailResponse is nil")
		}

		accountMap[grainId] = EmailResponse.Email
	}
	// grainIdsのlenとaccountMapのlenが一致しない場合はpanic
	if len(grainIds) != len(accountMap) {
		log.Printf("Error getting all account emails: %v", "grainIds and accountMap length mismatch")
		log.Printf("grainIds length: %v", len(grainIds))
		log.Printf("accountMap length: %v", len(accountMap))
		panic("grainIds and accountMap length mismatch")
	}
	return accountMap
}

func generateUUID() string {
	return uuid.New().String()
}
