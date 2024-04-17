# structured
GoでGraphQLを実装する最小構成に加え、よくある設定fileなどを分けたもの

## Overview
### 関係性
DB -- GraphQL server -- client

### GraphQL serverの責務
1. schema: mappingの定義
具体的には、以下の3つのmappingを担う
- GraphQLのqueryで使用されるfield名と型の定義
- GraphQLのfieldとapplicationのdata型(Goの構造体)のfiledのmapping
- GraphQLのfieldとDBのcolumnのmapping(必要な場合)
2. resolver: dataの取得
schemaによって定義されたmapping情報をもとに、実際にDBからdataを取得する
このとき、DAOを通じてDBのcolumnとapplicationのdata型(Goの構造体)をmappingすることもあれば、ORMで済ませられることもある

### GraphQL serverで実装すること
大きく2つ
- schemaの定義
- resolverの実装
特に、schemaの定義が重要。resolverの実装は一定形式的。

## Run
```
make up
$ curl -X POST -H "Content-Type: application/json" -d '{"query": "{ users { id name } }"}' http://localhost:8080/graphql
```