package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/cluster"
	c "github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/presentation"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func shutdownNodes(nodes []*cluster.Cluster) {
	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go func(n *cluster.Cluster) {
			defer wg.Done()
			n.Shutdown(true)
		}(node)
	}
	wg.Wait()
}

func startNodes(clusterName string, ports []int) []*cluster.Cluster {
	var nodes []*cluster.Cluster
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			node := c.StartNode(clusterName, port)
			nodes = append(nodes, node)
		}(port)
	}

	wg.Wait()
	return nodes
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	proto.ManagerFactory(func() proto.Manager {
		return &grain.ManagerGrain{}
	})

	// クラスタノードを起動
	ports := []int{6330, 6331}
	// ports := []int{6330}
	// ports := []int{6330, 6331, 6332, 6333, 6334}
	nodes := startNodes("my-cluster", ports)

	// nodeの安定化を待つ
	time.Sleep(10 * time.Second)

	// テストケース
	testCases := []int{1, 10, 100}
	// testCases := []int{100, 500, 1000, 5000, 10000}

	for _, totalRequests := range testCases {
		registeredIds, elapsed := registerAccounts(nodes, totalRequests)

		fmt.Printf("Total requests: %d, Elapsed time: %v\n", totalRequests, elapsed)

		// 処理がされきるまで一応待ってみる。待ったほうが成功率が上がるは上がる?わけでもないらしい
		time.Sleep(3 * time.Second)

		// registerは成功しているよう。IDがそれだけ生成されいている。
		// fmt.Printf("registeredIds: %v\n", registeredIds)

		// 登録されたアカウントの確認
		accounts := presentation.GetAllAccounts(nodes, registeredIds)
		// time.Sleep(10 * time.Second)

		fmt.Printf("Total Registered Emails: %d\n", len(accounts))
		for accountId, email := range accounts {
			fmt.Printf("%v : %v\n", accountId, email)
		}
	}
	// クラスタノードをシャットダウン
	shutdownNodes(nodes)
}

func registerAccounts(nodes []*cluster.Cluster, totalRequests int) ([]string, time.Duration) {
	var wg sync.WaitGroup
	requestRate := 10 // 1秒あたりの要求数
	requestInterval := time.Second / time.Duration(requestRate)
	registeredIds := make([]string, 0, totalRequests)

	start := time.Now()
	var mu sync.Mutex

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in goroutine %d: %v\n", i, r)
				}
				wg.Done()
			}()
			email := fmt.Sprintf("email%d@account.test", i)
			id, err := presentation.RegisterAccount(nodes[i%len(nodes)], email)
			if err != nil {
				log.Printf("Error registering account in goroutine %d: %v\n", i, err)
				return
			}
			if id == "" {
				log.Printf("Empty ID returned for email in goroutine %d: %s\n", i, email)
				return
			}
			mu.Lock()
			registeredIds = append(registeredIds, id)
			mu.Unlock()
			time.Sleep(requestInterval)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	return registeredIds, elapsed
}
