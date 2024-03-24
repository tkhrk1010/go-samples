# fruits-trading-demo
trade fruits with actor model and scatter-gather pattern

## Usecases
### Register Purchasing
```mermaid
graph TD
  Apple("Apple register API")
  Orange("Orange register API")
  Banana("Banana register API")
  InventoryDB("Inventory Management DB")

  Apple -->|"register"| InventoryDB
  Orange -->|"register"| InventoryDB
  Banana -->|"register"| InventoryDB
```

### Trade judgement

```mermaid
sequenceDiagram
  participant Trader
  participant TradeSupportInformationHandler 
  participant InventoryAggregator
  participant AppleInventoryCollector
  participant OrangeInventoryCollector
  participant BananaInventoryCollector
  participant MarketPriceCollector
  participant InventoryDB
  participant TradeResultDB

  Trader ->> TradeSupportInformationHandler : Start Transaction
  TradeSupportInformationHandler ->> MarketPriceCollector : Request Market Price
  TradeSupportInformationHandler ->> InventoryAggregator : Collect Information
  InventoryAggregator ->> AppleInventoryCollector : Get Inventory Information
  InventoryAggregator ->> OrangeInventoryCollector : Get Inventory Information
  InventoryAggregator ->> BananaInventoryCollector : Get Inventory Information
  AppleInventoryCollector ->> InventoryDB : Get Inventory Information
  InventoryDB -->> AppleInventoryCollector :  Inventory Information
  AppleInventoryCollector -->> InventoryAggregator : Inventory Information
  OrangeInventoryCollector ->> InventoryDB : Get Inventory Information
  InventoryDB -->> OrangeInventoryCollector :  Inventory Information
  OrangeInventoryCollector -->> InventoryAggregator : Inventory Information
  BananaInventoryCollector ->> InventoryDB : Get Inventory Information
  InventoryDB -->> BananaInventoryCollector :  Inventory Information
  BananaInventoryCollector -->> InventoryAggregator : Inventory Information
  InventoryAggregator -->> TradeSupportInformationHandler : Inventory Information
  MarketPriceCollector -->> TradeSupportInformationHandler : Market Price Information
  TradeSupportInformationHandler -->> Trader : Information
  Trader ->> Trader : judge sell/hold
  Trader ->> TradeResultDB : Save Result

```

## Actor structure
```mermaid
graph TD
  TradeSupportInformationHandler --> PriceCollector
  TradeSupportInformationHandler --> InventoryAggregator
  InventoryAggregator --> AppleInventoryCollector
  InventoryAggregator --> OrangeInventoryCollector
  InventoryAggregator --> BananaInventoryCollector
```
## QuickStart
```
$ docker-compose up -d
$ go run .
```