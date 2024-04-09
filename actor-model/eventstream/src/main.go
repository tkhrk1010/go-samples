package main

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
)

// event streamに発行されるMessageの型
type Message struct {
	Text string
}

// event streamに発行されるイベントの型
type ESEvent struct {
	Message Message
}

// PublisherActorはMessageを受け取り、event streamにpublishする
type PublisherActor struct {
	stream *eventstream.EventStream
}

func NewPublisherActor(stream *eventstream.EventStream) *PublisherActor {
	return &PublisherActor{
		stream: stream,
	}
}

func (state *PublisherActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
  // PublisherActorがMessageを受け取ると、event streamにpublishする
	case Message:
		event := ESEvent{Message: msg}
		state.stream.Publish(event)
	}
}

// 使ってない。SubscriberActorがstreamをSubscribeするためにはどうする？
// あくまでevent streamにもPIDがあるので、ESがactorともいえる。はず。
// type SubscriberActor struct{}

func SubscriberReceive(context actor.Context) {
	switch msg := context.Message().(type) {
	case ESEvent:
		fmt.Printf("Subscriber received message: %v\n", msg.Message.Text)
	}
}

func main() {
	fmt.Println("start")
	system := actor.NewActorSystem()
  stream := eventstream.NewEventStream()

	fmt.Println("create publisherActor")
	props := actor.PropsFromProducer(func() actor.Actor { return NewPublisherActor(stream) })
	publisherActorPid, _ := system.Root.SpawnNamed(props, "publisherActor")

	fmt.Println("subscibing...")
	stream.Subscribe(func(event interface{}) {
		system.Root.Send(system.Root.Spawn(actor.PropsFromFunc(SubscriberReceive)), event)
	})

	fmt.Println("send message to publisherActor")
	message := Message{Text: "Hello, Protoactor!"}
	system.Root.Send(publisherActorPid, message)

	fmt.Println("wait for message to be published")
  // goroutineを終了しないようにする
	console := make(chan struct{})
  // channelからの受信を待つ
	<-console
}
