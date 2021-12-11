
# get abigen:
# cd $GOPATH/pkg/github.com/ethereum/go-ethereum/
# make devtools

abigen --abi=ethBridge.json --pkg=contracts --type=eth --out=ethBridge.go
abigen --abi=ambBridge.json --pkg=contracts --type=amb --out=ambBridge.go
