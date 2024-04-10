package main

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
)

type ParentActor struct {
	childPID *actor.PID
}

func NewParentActor() actor.Actor {
	return &ParentActor{}
}

func (state *ParentActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Println("Parent actor started")
		log.Println("Parent actor messege typed *actor.Started: ", msg)

	case CreateChild:
		log.Println("Parent actor received CreateChild message")
		// ChildActorを作成
		state.childPID = context.Spawn(actor.PropsFromProducer(func() actor.Actor { return &ChildActor{} }))
		// ParentActorはChildActorを監視する
		context.Watch(state.childPID)

	case MakeChildPanic:
		log.Println("Parent actor received MakeChildPanic message")
		// 子アクターにpanicを発生させる
		// ここにchildPIDを指定して、親アクターから子アクターにメッセージを送信する
		context.Send(state.childPID, MakeChildPanic{})
	}

}

type CreateChild struct{}

type MakeChildPanic struct{}

type ChildActor struct{}

func (state *ChildActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Println("Child actor started")
		log.Println("Child actor messege typed *actor.Started: ", msg)
		log.Println("Child actor's Address: ", context.Self().Address, ", PID: ", context.Self().Id)
	case MakeChildPanic:
		log.Println("Child actor received MakeChildPanic message")
		// panicを発生させる
		panic("Something very bad happened")

	// 自分(子アクター)が停止した場合、親アクターに通知したりできる。
	// ここでは、単にログを出力している
	case *actor.Stopped:
		log.Printf("Child actor stopped. ID: %v\n", context.Self().Id)
	}
	
}

func main() {
	// microsecondまでtimestampを表示する
	log.SetFlags(log.Lmicroseconds)

	log.Println("start")
	// ActorSystemを作成
	system := actor.NewActorSystem()

	//
	// ParentActorを作成
	log.Println("create ParentActor")
	// 監視戦略の定義
	// ここでは、子アクターがpanicしたらその子アクターを停止する
	decider := func(reason interface{}) actor.Directive {
		// panicの理由を出力
		log.Printf("Child actor panicked: %v\n", reason)
		// 子アクターを再起動せず、停止する
		return actor.StopDirective
	}
	// supervisor strategyを設定
	// NewOneForOneStrategyは子アクターごとに個別の監視を行うための戦略
	supervisor := actor.NewOneForOneStrategy(1, 1, decider)
	// https://pkg.go.dev/github.com/asynkron/protoactor-go@v0.0.0-20240408180828-2a5372b5b465/actor#PropsOption
	// ()をつけずに渡すことで、関数を呼び出すのではなく、関数自体を渡すことができる
	// これにより、PropsFromProducerは必要なタイミングで自由に関数を呼び出すことができる
	parentProps := actor.PropsFromProducer(NewParentActor, actor.WithSupervisor(supervisor))
	parentPID := system.Root.Spawn(parentProps)
	log.Printf("Parent actor PID: %v\n", parentPID)

	// ParentActorにCreateChildメッセージを送信し、子アクターを作成してもらう
	log.Println("send CreateChild message to ParentActor")
	system.Root.Send(parentPID, CreateChild{})

	// 子アクターにpanicを発生させる
	// system.Rootは子アクターのことを知らないので、親アクター経由でmessageを伝達するのが一般的
	// こういう構成にすることで、子アクターの存在はカプセル化されて、親アクターが子アクターの管理を行いやすい
	log.Println("send MakeChildPanic message to ParentActor")
	system.Root.Send(parentPID, MakeChildPanic{})

	//
	// main関数が終了しないようにする一般的な方法
	// 大元のmain関数をすぐ終了させてしまうと、非同期でmessageを受け取ったactorが処理の途中で終了してしまうため
	console := make(chan struct{})
	// console channelからのdataを受信する機能だが、実際に受信することはないので、関数が終了しない
	<-console
}
