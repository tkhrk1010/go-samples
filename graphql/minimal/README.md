# minimal
GoでGraphQLを実装する最小構成


## Run
```
make up
curl -X POST -H "Content-Type: application/json" -d '{"query": "{ users { id name } }"}' http://localhost:8080/graphql
```