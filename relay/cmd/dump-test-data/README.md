# Dump Test Data

## Usage:

go run cmd/dump-test-data/main.go [comamnd] [arg]

### Comamnds:
+ receipts-proof [network tag]
+ pow-block [block number]
+ poa-block [block number]
+ epoch [epoch number]

Network tags:
+ amb
+ eth

## Examples:
Generating ambrosus testing data for receipts proof:
```sh
go run cmd/dump-test-data/main.go receipts-proof amb
```

Encoding ethereum PoW block:
```sh
go run cmd/dump-test-data/main.go pow-block 14257704
```