## 実装するもの
persistence/persistence_provider.go
```go
// Provider is the abstraction used for persistence
type Provider interface {
	GetState() ProviderState
}

// ProviderState is an object containing the implementation for the provider
type ProviderState interface {
	SnapshotStore
	EventStore

	Restart()
	GetSnapshotInterval() int
}

type SnapshotStore interface {
	GetSnapshot(actorName string) (snapshot interface{}, eventIndex int, ok bool)
	PersistSnapshot(actorName string, snapshotIndex int, snapshot proto.Message)
	DeleteSnapshots(actorName string, inclusiveToIndex int)
}

type EventStore interface {
	GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{}))
	PersistEvent(actorName string, eventIndex int, event proto.Message)
	DeleteEvents(actorName string, inclusiveToIndex int)
}
```

## 設計
- [x] methodの洗い出し

## must実装
- [ ] ProviderState(DynamoDBProvider)の実装
  - [ ] ProviderState構造体の実装
  - [ ] GetStateの実装
- [ ] SnapshotStoreの実装
  - [ ] SnapshotStore構造体(entry)の実装
  - [ ] GetSnapshotの実装
  - [ ] PersistSnapshotの実装
  - [ ] DeleteSnapshotsの実装
- [ ] EventStoreの実装
  - [ ] EventStore構造体(entry)の実装
  - [ ] GetEventsの実装
  - [ ] PersistEventの実装
  - [ ] DeleteEventsの実装
- [ ] Restartの実装(空method)
- [ ] GetSnapshotIntervalの実装
# more実装
- [ ] NewDynamoDBProviderの実装