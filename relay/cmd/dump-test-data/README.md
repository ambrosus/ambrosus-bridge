# Dump Test Data

## Usage:

go run cmd/dump-test-data/main.go [comamnd] [arg]

### Comamnds:
+ receipts-proof [network]
+ pow-block
+ poa-block

Networks:
+ amb
+ eth

## Examples:
Generating ambrosus testing data for receipts proof:
```sh
go run cmd/dump-test-data/main.go receipts-proof amb
```
Encoding PoW block.
```sh
go run cmd/dump-test-data/main.go pow-block
```