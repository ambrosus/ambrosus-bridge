# Download abigen:
# https://geth.ethereum.org/downloads/


yarn --cwd ../../contracts --silent ts-node ./scripts/merge_abi.ts --relay | abigen --abi - --pkg=bindings --type=bridge --out=../internal/bindings/bridge.go
abigen --abi=../../contracts/abi/ERC20.json --pkg=bindings --type=token --out=../internal/bindings/erc20.go
