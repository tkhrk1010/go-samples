package main

import (
	"fmt"
	"time"

	"github.com/Workiva/go-datastructures/queue"
)

type Message struct {
	text string
}

type Actor struct {
	mailbox *queue.Queue
	stopCh  chan bool
}

func NewActor() *Actor {
	return &Actor{
		mailbox: queue.New(100),
		stopCh:  make(chan bool),
	}
}

func (a *Actor) Start() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in Actor, ", r)
			}
		}()
		for {
			select {
			case <-a.stopCh:
				return
			default:
				message, _ := a.mailbox.Get(1)
				if message != nil {
					fmt.Println(message[0].(Message).text)
					if message[0].(Message).text == "error" {
						panic("Something bad happened")
					}
				}
			}
		}
	}()
}

func (a *Actor) Stop() {
	close(a.stopCh)
}

func (a *Actor) Send(msg Message) {
	a.mailbox.Put(msg)
}

func main() {
	actor1 := NewActor()
	actor1.Start()
	defer actor1.Stop()

	actor2 := NewActor()
	actor2.Start()
	defer actor2.Stop()

	actor1.Send(Message{text: "Hello from Actor 1"})
	actor2.Send(Message{text: "Hello from Actor 2"})
	actor1.Send(Message{text: "error"})  // This will cause actor1 to panic

	time.Sleep(2 * time.Second)  // Give some time for the messages to be processed
}
