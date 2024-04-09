# Go actor model
Actor modelの実装

# Why actor model?
各Actorは独立したgoroutineで動作し、独自のメッセージキューを持っているため、一つのActorが失敗しても他のActorに影響を与えない

