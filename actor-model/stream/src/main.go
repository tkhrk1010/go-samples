// 
package main

import (
    "fmt"
    "time"

    "github.com/asynkron/protoactor-go/actor"
    "github.com/asynkron/protoactor-go/stream"
)

// Define the Message struct
type Message struct {
    Text string
}

func main() {
    fmt.Println("start")

    // Create a new ActorSystem
    system := actor.NewActorSystem()

    // Create a new typed stream for Message
		// Message型のstreamを作成
    s := stream.NewTypedStream[Message](system)

    // Spawn a goroutine to send messages to the stream
    go func() {
        for i := 0; i < 5; i++ {
            msg := Message{Text: fmt.Sprintf("Message %d", i)}
            system.Root.Send(s.PID(), msg)
						// 少し間隔を開けて送信
            time.Sleep(100 * time.Millisecond)
        }
				// 5回送ったら、streamを閉じる
        s.Close()
    }()

    // Receive messages from the stream and print them
		// .C()はchannelで、streamから受け取ったメッセージを受け取る
		// range s.C()とすることで、streamがcloseされたらループが終了する
    for msg := range s.C() {
        fmt.Printf("Received: %v\n", msg.Text)
    }

    fmt.Println("end")
}