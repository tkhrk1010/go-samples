// protoactorのclusterにnodeを登録するための関数を提供する
package cluster

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

func StartNode(port int) *cluster.Cluster {
	system := actor.NewActorSystem()

	provider := automanaged.NewWithConfig(2*time.Second, port, "localhost:6330", "localhost:6331")
	lookup := disthash.New()
	config := remote.Configure("localhost", 0)

	accountKind := proto.NewAccountKind(func() proto.Account {
		return &grain.AccountGrain{}
	}, 0)

	managerKind := proto.NewManagerKind(func() proto.Manager {
		return &grain.ManagerGrain{}
	}, 0)

	clusterConfig := cluster.Configure("my-cluster", provider, lookup, config,
		cluster.WithKinds(accountKind, managerKind))

	cluster := cluster.New(system, clusterConfig)

	cluster.StartMember()

	return cluster
}
