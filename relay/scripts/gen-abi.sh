# Download abigen:
# https://geth.ethereum.org/downloads/

abigen --abi=../../contracts/abi/EthBridge.json --pkg=contracts --type=eth --out=../contracts/ethBridge.go
abigen --abi=../../contracts/abi/AmbBridge.json --pkg=contracts --type=amb --out=../contracts/ambBridge.go
abigen --abi=../../contracts/abi/ValidatorSet.json --pkg=contracts --type=vs --out=../contracts/validatorSet.go
