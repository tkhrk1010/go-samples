// protoactorのclusterにnodeを登録するための関数を提供する
package cluster

import (
	"time"
	"fmt"
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

// port割当
func getUsableLocalPorts() []string {
	var addresses []string

	for port := 6330; port <= 6336; port++ {
			address := fmt.Sprintf("localhost:%d", port)
			listener, err := net.Listen("tcp", address)
			if err != nil {
					continue
			}
			listener.Close()
			addresses = append(addresses, address)
	}

	return addresses
}

func StartNode(clasterName string, port int) *cluster.Cluster {
	system := actor.NewActorSystem()

	provider := automanaged.NewWithConfig(2*time.Second, port, getUsableLocalPorts()...)
	lookup := disthash.New()
	config := remote.Configure("localhost", 0)

	managerKind := proto.NewManagerKind(func() proto.Manager {
		return &grain.ManagerGrain{}
	}, 0)

	clusterConfig := cluster.Configure(clasterName, provider, lookup, config,
		cluster.WithKinds(managerKind))

	cluster := cluster.New(system, clusterConfig)

	cluster.StartMember()

	return cluster
}
