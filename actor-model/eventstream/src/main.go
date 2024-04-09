package main

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/eventstream"
)

// event streamをsubscribeするActor
type SubscriberActor struct{}

func (state *SubscriberActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case ESEvent:
		fmt.Printf("Subscriber received message: %v\n", msg.Message.Text)
	}
}

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


func main() {
	fmt.Println("start")
	system := actor.NewActorSystem()
	stream := eventstream.NewEventStream()

	fmt.Println("create subscriberActor")
  subscriberProps := actor.PropsFromProducer(func() actor.Actor { return &SubscriberActor{} })
  subscriberPID := system.Root.Spawn(subscriberProps)

  fmt.Println("subscribe to event stream")
  stream.Subscribe(func(event interface{}) {
      system.Root.Send(subscriberPID, event)
  })

  fmt.Println("create publisherActor")
	publisherProps := actor.PropsFromProducer(func() actor.Actor { return NewPublisherActor(stream) })
	publisherPid, _ := system.Root.SpawnNamed(publisherProps, "publisherActor")

	fmt.Println("send message to publisherActor")
	message := Message{Text: "Hello, Protoactor!"}
	system.Root.Send(publisherPid, message)

	fmt.Println("wait for message to be published")
	// goroutineを終了しないようにする
	console := make(chan struct{})
	// channelからの受信を待つ
	<-console
}
