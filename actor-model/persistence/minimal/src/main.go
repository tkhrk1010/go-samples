package main

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
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

	Value string
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
		se := &event.MyEvent{Value: a.Value}
		a.PersistSnapshot(se)
		log.Printf("Snapshot persisted: %v", se)
	case *persistence.ReplayComplete:
		log.Println("ReplayComplete received")
	case *event.MyEvent:
		log.Printf("Received MyEvent: %v", msg.Value)
		a.Value = msg.Value
		// 永続化
		a.PersistReceive(msg)
	default:
		log.Printf("Unknown message received: %v", msg)
	}
}

func getSnapshot(provider *MyInmemoryProvider, myActorPid *actor.PID) {
	// スナップショットを取得
	snapshot, eventIndex, ok := provider.providerState.GetSnapshot(myActorPid.Id)
	if ok {
		log.Printf("Retrieved snapshot: %v, eventIndex: %d", snapshot, eventIndex)
	} else {
		log.Println("Snapshot not found")
	}
}

func main() {
	log.Println("start")
	system := actor.NewActorSystem()

	//
	// 永続化のprovider(persistence)を作成
	log.Println("create provider")
	// 引数はsnapshotを取るinterval
	provider := NewMyInmemoryProvider(2)

	//
	// MyActorを生成
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
	myActorPid := system.Root.Spawn(props)
	log.Printf("MyActor PID: %s", myActorPid)

	//
	// eventを永続化
	log.Println("persist event")
	// arg: actorName string, eventIndex int, event proto.Message
	// https://github.com/asynkron/protoactor-go/blob/2a5372b5b465b3bb030dd26086cb5840465e7354/persistence/in_memory_provider.go#L79
	// InmemoryProviderにeventを渡すときは、proto.Messageを渡す必要がある
	// そのため、eventを作っている
	getSnapshot(provider, myActorPid)
	system.Root.Send(myActorPid, &event.MyEvent{Value: "first message: please sum =+ 1"})
	time.Sleep(1 * time.Second)
	getSnapshot(provider, myActorPid)
	system.Root.Send(myActorPid, &event.MyEvent{Value: "second message: please sum =+ 1"})
	time.Sleep(1 * time.Second)
	getSnapshot(provider, myActorPid)
	system.Root.Send(myActorPid, &event.MyEvent{Value: "third message: please sum =+ 1"})
	time.Sleep(1 * time.Second)
	getSnapshot(provider, myActorPid)
	system.Root.Send(myActorPid, &event.MyEvent{Value: "fourth message: please sum =+ 1"})
	time.Sleep(1 * time.Second)
	getSnapshot(provider, myActorPid)
	system.Root.Send(myActorPid, &event.MyEvent{Value: "fifth message: please sum =+ 1"})

	// actorが永続化してくれるのを少し待つ
	time.Sleep(2 * time.Second)

	//
	// eventを取得
	log.Println("get event")
	// actorNameは、context.Self().Idを使っているらしいことが公式からわかる。
	// https://github.com/asynkron/protoactor-go/blob/2a5372b5b465b3bb030dd26086cb5840465e7354/persistence/plugin.go#L55
	// そのため、actorNameは、actorのPIDを使う
	// 引数2: eventIndexStartは、取得したい最初のindexを指定する。
	// 引数3: eventIndexEndは、取得したい最後のindexを指定する。全体のindexを超えているとerrorになる。0を指定すると全て取得らしい
	// https://github.com/asynkron/protoactor-go/blob/2a5372b5b465b3bb030dd26086cb5840465e7354/persistence/in_memory_provider.go#L69C84-L69C97
	provider.GetEvents(myActorPid.Id, 0, 0)

	// systemが終了しないように待機
	console := make(chan string)
	<-console
}
