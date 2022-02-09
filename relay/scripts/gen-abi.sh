# Download abigen:
# https://geth.ethereum.org/downloads/

abigen --abi=./contracts/ethBridge.json --pkg=contracts --type=eth --out=./contracts/ethBridge.go
abigen --abi=./contracts/ambBridge.json --pkg=contracts --type=amb --out=./contracts/ambBridge.go
