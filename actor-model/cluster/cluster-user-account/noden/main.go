// nodeがn個ないとhealcheckでerrorが出続けてうるさいので、全node立ち上げる
package main

import (
	"github.com/asynkron/goconsole"
	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/cluster"
)

func main() {
	cluster.StartNode("my-cluster", 6332)
	cluster.StartNode("my-cluster", 6333)
	cluster.StartNode("my-cluster", 6334)
	_, _ = console.ReadLine()
}