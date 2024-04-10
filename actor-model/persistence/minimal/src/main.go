package main

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	// "google.golang.org/protobuf/proto"
	"github.com/tkhrk1010/go-samples/actor-model/persistence/minimal/src/event"
)

// ref:
// https://github.com/Huawei-CPLLab/protoactor-go/blob/1579a3e2bb055995f842e2a70284bc37b8e7e545/examples/persistence/main.go#L9
// https://github.com/asynkron/protoactor-go/blob/2a5372b5b465b3bb030dd26086cb5840465e7354/persistence/persistence_provider.go
type MyInmemoryProvider struct {
	providerState persistence.ProviderState
}

func NewMyInmemoryProvider(snapshotInterval int) *MyInmemoryProvider {
	return &MyInmemoryProvider{
		providerState: persistence.NewInMemoryProvider(snapshotInterval),
	}
}

func (p *MyInmemoryProvider) GetState() persistence.ProviderState {
	return p.providerState
}

func (p *MyInmemoryProvider) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int) {
	callback := func(e interface{}) {
		if msg, ok := e.(*event.MyEvent); ok {
			log.Printf("Retrieved event: %v", msg.Value)
		} else {
			log.Printf("Unknown event type: %T", e)
		}
	}
	p.providerState.GetEvents(actorName, eventIndexStart, eventIndexEnd, callback)
}

type MyActor struct {
	// このように構造体を書くことで、定義した構造体のfieldやmethodをすべて使えるようになる
	// persistence.Mixinはtypeであり、actorの永続化に必要なものが色々定義されている。
	// これを継承しないと、永続化可能なactorとして機能しない
	persistence.Mixin
}

func NewMyActor() actor.Actor {
	return &MyActor{}
}

func (a *MyActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Println("MyActor started")
	case *persistence.RequestSnapshot:
		log.Println("RequestSnapshot received")
	case *persistence.ReplayComplete:
		log.Println("ReplayComplete received")
	default:
		log.Printf("Unknown message received: %v", msg)
	}
}


func main() {
	log.Println("start")
	system := actor.NewActorSystem()

	log.Println("create provider")
	// 引数はsnapshotを取るinterval
	provider := NewMyInmemoryProvider(3)

	log.Println("create MyActor")
	// Providor型である必要がある？というよりは、providerStateというfiledを持っている必要がある？
	// Usingを使ってactor.ReceiverFuncを定義する。これが永続化の初期化のやり方らしい
	// 公式:
	// https://github.com/asynkron/protoactor-go/blob/2a5372b5b465b3bb030dd26086cb5840465e7354/persistence/receiver.go
	// このとき、引数はProvider interfaceの実装型である必要があり、GetState()を持っている必要がある
	// GetStateはProviderStateを返す
	// providerStateとは、永続化の実装をを持つobjectのこと。以下公式
	// // ProviderState is an object containing the implementation for the provider
	// from https://github.com/asynkron/protoactor-go/blob/2a5372b5b465b3bb030dd26086cb5840465e7354/persistence/persistence_provider.go#L8
	// 実装の参考:
	// https://github.com/ytake/protoactor-go-persistence-pg
	// https://github.com/asynkron/protoactor-go/blob/2a5372b5b465b3bb030dd26086cb5840465e7354/persistence/plugin_test.go
	props := actor.PropsFromProducer(NewMyActor, actor.WithReceiverMiddleware(persistence.Using(provider)))

	// prosを使ってactorを生成
	pid := system.Root.Spawn(props)
	log.Printf("Actor PID: %s", pid)


	log.Println("persist event")
	// eventを永続化
	// arg: actorName string, eventIndex int, event proto.Message
	// https://github.com/asynkron/protoactor-go/blob/2a5372b5b465b3bb030dd26086cb5840465e7354/persistence/in_memory_provider.go#L79
	// InmemoryProviderにeventを渡すときは、proto.Messageを渡す必要がある
	// そのため、eventを作っている
	provider.providerState.PersistEvent("testActor", 0, &event.MyEvent{Value: "first message"})
	provider.providerState.PersistEvent("testActor", 1, &event.MyEvent{Value: "second message"})

	// eventを取得
	log.Println("get event")
	provider.GetEvents("testActor", 0, 2)

	// systemが終了しないように数秒待つ
	time.Sleep(3)
}
