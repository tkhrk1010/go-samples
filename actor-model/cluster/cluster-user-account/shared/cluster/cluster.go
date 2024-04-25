// protoactorのclusterにnodeを登録するための関数を提供する
package cluster

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/lmittmann/tint"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/grain"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/proto"
)

// port割当
func getUsableLocalPorts() []string {
	var addresses []string

	for port := 6330; port <= 6334; port++ {
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

// logger
// ref: /protoactor-go/examples/actor-logging/main.go
func coloredConsoleLogging(system *actor.ActorSystem) *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelWarn,
		TimeFormat: time.RFC3339,
		AddSource:  true,
	})).With("lib", "Proto.Actor").
		With("system", system.ID)
}

func StartNode(clasterName string, port int) *cluster.Cluster {
	system := actor.NewActorSystem(actor.WithLoggerFactory(coloredConsoleLogging))

	provider := automanaged.NewWithConfig(2*time.Second, port, getUsableLocalPorts()...)
	lookup := disthash.New()
	config := remote.Configure("localhost", 0)

	AmccountKind := proto.NewAccountKind(func() proto.Account {
		return &grain.AccountGrain{}
	}, 0)

	clusterConfig := cluster.Configure(clasterName, provider, lookup, config,
		cluster.WithKinds(AmccountKind))

	cluster := cluster.New(system, clusterConfig)

	cluster.StartMember()

	return cluster
}
