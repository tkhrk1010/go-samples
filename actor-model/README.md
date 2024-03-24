# Go actor model Sample
Actor modelを試しに実行してみるためのsample

# Quick start
```
$ go mod init github.com/your_name/your_repository_name
$ go get github.com/asynkron/protoactor-go
$ go mod tidy
$ go run .
```

# Why actor model?
各Actorは独立したgoroutineで動作し、独自のメッセージキューを持っているため、一つのActorが失敗しても他のActorに影響を与えない

