# event-sourcing

風速をeventとしてdynamoDBに保存する
→ dynamodb streamsが受信
→ read model updaterがRDBを更新

以下を使用  
https://github.com/tkhrk1010/protoactor-go-persistence-dynamodb

rmuは、以下を使用させていただきました  
https://github.com/j5ik2o/cqrs-es-example-go/tree/main/pkg/rmu


TODO:
- [x] Event, Snapshotにtypeとidを追加
- [ ] rmu呼び出しロジック
