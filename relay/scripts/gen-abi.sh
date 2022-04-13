# Download abigen:
# https://geth.ethereum.org/downloads/

# abigen --abi=../../contracts/abi/EthBridge.json --pkg=contracts --type=eth --out=../internal/contracts/ethBridge.go
# abigen --abi=../../contracts/abi/AmbBridge.json --pkg=contracts --type=amb --out=../internal/contracts/ambBridge.go
python merge_abis.py | abigen --abi - --pkg=contracts --type=common --out=../internal/contracts/common.go
abigen --abi=../../contracts/abi/ModifiedValidatorSet.json --pkg=contracts --type=vs --out=../internal/contracts/validatorSet.go
