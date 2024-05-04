# event-sourcing

風速をeventとしてdynamoDBに保存する
→ dynamodb streamsが受信
→ read model updaterがRDBを更新

以下を使用  
https://github.com/tkhrk1010/protoactor-go-persistence-dynamodb