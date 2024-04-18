# k8s-first
proto-actor-goをk8s上で動かすための最小構成のつもりで始めた  
作ってて、最小構成じゃない感じにはなってきた  
正直あんまり意味のある最小がわからない。

## Quick Start
```
$ minikube start
$ make run
```

## 構築手順
```
$ cd parent
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/actor.proto
$ cd ../child
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/actor.proto
```
