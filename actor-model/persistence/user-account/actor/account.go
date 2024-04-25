package actor

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/tkhrk1010/go-samples/actor-model/persistence/user-account/proto"
)

type AccountActor struct {
	persistence.Mixin
	// id 		string
	email 	string
}


func NewSignUpEvent(email string) *proto.AccountEvent {
	return &proto.AccountEvent{
			Metadata: &proto.EventMetadata{
					Type:       "SignUp",
					OccurredAt: time.Now().Format(time.RFC3339),
			},
			Content: &proto.AccountEvent_SignUp{
					SignUp: &proto.SignUp{
							Email: email,
					},
			},
	}
}

func NewSignUpSnapshot(email string) *proto.AccountSnapshot {
	return &proto.AccountSnapshot{
			Metadata: &proto.EventMetadata{
					Type:       "SignUp",
					OccurredAt: time.Now().Format(time.RFC3339),
			},
			Content: &proto.AccountSnapshot_SignUp{
					SignUp: &proto.SignUp{
							Email: email,
					},
			},
	}
}

func (a *AccountActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Println("actor started")
	case *persistence.RequestSnapshot:
		log.Printf("snapshot actor email '%v'", a.email)
		a.PersistSnapshot(NewSignUpSnapshot(a.email))
	case *proto.AccountSnapshot:
		switch content := msg.Content.(type) {
    case *proto.AccountSnapshot_SignUp:
      a.email = content.SignUp.Email
			log.Printf("recovered from snapshot, actor email changed to '%v'", a.email)
    case *proto.AccountSnapshot_Login:
			// implement
    case *proto.AccountSnapshot_Logout:
			// implement
		}
	case *persistence.ReplayComplete:
		log.Printf("replay completed, actor email changed to '%v'", a.email)
	case *proto.AccountEvent:
		scenario := "received replayed event"
		if !a.Recovering() {
			a.PersistReceive(msg)
			scenario = "received new message"
		}
		switch content := msg.Content.(type) {
		case *proto.AccountEvent_SignUp:
			a.email = content.SignUp.Email
			log.Printf("%s, actor email changed to '%v'\n", scenario, a.email)
		case *proto.AccountEvent_Login:
			// implement
    case *proto.AccountEvent_Logout:
			// implement
		}
	}
}
