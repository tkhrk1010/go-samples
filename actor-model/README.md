# Go actor model
Actor modelの実装

## Why actor model?
各Actorは独立したgoroutineで動作し、独自のメッセージキューを持っているため、一つのActorが失敗しても他のActorに影響を与えない

## 公式example
https://github.com/asynkron/protoactor-go/tree/76c172a71a1687fbf31d45dd6ec33e005dfd6eaf/examples  
まずは自分でつくらなくなても、これでいい、、

## claude3のよくある間違いメモ
- ReceiverFuncを使う → actor.Behavior.Receiveを使ってもらう
- .WithReceiverMiddlewareを使う → props := actor.PropsFromProducer(NewMyActor, actor.WithReceiverMiddleware())のようにして、propsを作ってもらう
- PropsFromFuncを使う →PropsFromProducerを使ってもらう 
- .Spawnを使う, ProduceActorを使う → system := actor.NewActorSystem()して、system.Root.Spawn(workerProps)してもらう。子アクターを作る場合は、actor.Context.Spawnしてもらう。
- NewRoundRobinPoolに.WithProducerを使う → pid := rootContext.Spawn(router.NewRoundRobinPool(5, actor.WithFunc(func)))を使ってもらう。
https://github.com/asynkron/protoactor-go/blob/76c172a71a1687fbf31d45dd6ec33e005dfd6eaf/examples/router-demo/main.go#L31
- .Tellを使う → .Sendを使う