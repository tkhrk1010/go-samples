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
- [ ] Restartの実装(しない)
- [ ] GetSnapshotIntervalの実装
# more実装
- [ ] NewDynamoDBProviderの実装