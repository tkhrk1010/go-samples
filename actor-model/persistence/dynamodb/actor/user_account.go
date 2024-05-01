package actor

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	p "github.com/tkhrk1010/go-samples/actor-model/persistence/dynamodb/persistence"
)

type UserAccount struct {
	persistence.Mixin
	Name  string
	Email string
}

func (u *UserAccount) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *p.Event:
		u.Name = msg.GetData()
		log.Printf("User account name updated to: %s", u.Name)
	case *p.Snapshot:
		// 特に何もしない
	case *persistence.RequestSnapshot:
		// PersistSnapshotを呼び出してスナップショットを保存する
		u.PersistSnapshot(ctx)
	default:
		log.Printf("Unknown message: %v", msg)
	}
}

func (u *UserAccount) PersistSnapshot(ctx actor.Context) {
	if u.Name == "" {
		// 空の状態はスナップショットを保存しない
		return
	}
	snapshot := &p.Snapshot{
		Data: u.Name,
	}
	u.Mixin.PersistSnapshot(snapshot)
}

func NewUserAccount() actor.Actor {
	return &UserAccount{}
}