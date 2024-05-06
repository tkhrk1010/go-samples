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
- [x] rmu呼び出しロジック
- [ ] time=2024-05-06T14:14:30.409Z level=INFO msg="{ "errorMessage ": "illegal base64 data at input byte 1 ", "errorType ": "CorruptInputError "}"

