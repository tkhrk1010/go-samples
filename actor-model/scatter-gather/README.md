# go-scatter-gather-sample
protoactor-go scatter-gather pattern sample code

## samples
- [simple](./simple): simple app
- [db](./db): save data in db
- [fruits-trading-demo](./fruits-trading-demo): trading demo app

## Structure of simple/ and db/
### model
```mermaid
classDiagram
    class System {
      +Aggregate Data
      +Collect Temperature Data
      +Collect Humidity Data
      +Collect Wind Speed Data
    }

    class User {
      +Request Data Aggregation
    }

    User --> System : Uses
    System --> AggregateData : Uses
    System --> CollectTemperatureData : Uses
    System --> CollectHumidityData : Uses
    System --> CollectWindSpeedData : Uses

    class AggregateData {
      +Integrate Collected Data
    }

    class CollectTemperatureData {
      +Gather Temperature Data
    }

    class CollectHumidityData {
      +Gather Humidity Data
    }

    class CollectWindSpeedData {
      +Gather Wind Speed Data
    }
```

### Actor tree
```mermaid
graph TD
    Main["main.go (Root Actor)"]
    Aggregator["AggregatorActor"]
    TempCollector["TemperatureCollectorActor"]
    HumidCollector["HumidityCollectorActor"]
    WindCollector["WindSpeedCollectorActor"]

    Main --> Aggregator
    Aggregator --> TempCollector
    Aggregator --> HumidCollector
    Aggregator --> WindCollector
```

### sequence
```mermaid
sequenceDiagram
    participant Main as main.go (Root Actor)
    participant Aggregator as AggregatorActor
    participant TempCollector as TemperatureCollectorActor
    participant HumidCollector as HumidityCollectorActor
    participant WindCollector as WindSpeedCollectorActor

    Main->>Aggregator: RequestFuture (AggregateRequest)
    Aggregator->>TempCollector: Request (CollectFeatureRequest)
    Aggregator->>HumidCollector: Request (CollectFeatureRequest)
    Aggregator->>WindCollector: Request (CollectFeatureRequest)
    TempCollector->>Aggregator: Respond (FeatureResponse)
    HumidCollector->>Aggregator: Respond (FeatureResponse)
    WindCollector->>Aggregator: Respond (FeatureResponse)
    Aggregator->>Main: Respond (AggregateResponse)

```


# Run
### simple
```
$ cd simple
$ go run .
```

### save result to postgresql db
```
$ cd db
$ docker-compose up -d
$ go run .
```
and manually check the db data

# init go
```
$ go mod init github.com/your_github_username/your_repository_name
$ go get github.com/asynkron/protoactor-go
$ go mod tidy
$ go run .
```
