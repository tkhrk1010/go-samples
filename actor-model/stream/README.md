proto-actor-goをつかって、actor modelにおけるstreamの使い方をざっくりと学ぶ

https://github.com/asynkron/protoactor-go/tree/2a5372b5b465b3bb030dd26086cb5840465e7354/stream


EventSourcingのstreamは別にeventstreamsがある  
https://github.com/asynkron/protoactor-go/tree/2a5372b5b465b3bb030dd26086cb5840465e7354/eventstream  

streaming処理にはstream、永続化をしたりするevent soucingにはeventstreamがいいらしい。

Akka実践バイブルの12章に、stream処理の概要が書いてある。