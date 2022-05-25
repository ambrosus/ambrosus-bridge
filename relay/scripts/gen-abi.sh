# Download abigen:
# https://geth.ethereum.org/downloads/

# abigen --abi=../../contracts/abi/EthBridge.json --pkg=contracts --type=eth --out=../internal/contracts/ethBridge.go
# abigen --abi=../../contracts/abi/AmbBridge.json --pkg=contracts --type=amb --out=../internal/contracts/ambBridge.go
python merge_abis.py ETH_AmbBridge.json ETH_EthBridge.json BSC_AmbBridge.json BSC_BscBridge.json | abigen --abi - --pkg=contracts --type=bridge --out=../internal/contracts/bridge.go
abigen --abi=../../contracts/abi/ValidatorSet.json --pkg=contracts --type=vs --out=../internal/contracts/validatorSet.go
abigen --abi=../../contracts/abi/IERC20.json --pkg=contracts --type=token --out=../internal/contracts/erc20.go
