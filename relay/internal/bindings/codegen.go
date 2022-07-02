package bindings

// Need to install "interfacer" and "mockgen"
// go install github.com/rjeczalik/interfaces/cmd/interfacer@latest
// go install github.com/golang/mock/mockgen@v1.6.0

//go:generate interfacer -for github.com/ambrosus/ambrosus-bridge/relay/internal/bindings.Bridge -as interfaces.BridgeContract -o ./interfaces/bridge.go
//go:generate mockgen -source=./interfaces/bridge.go -destination=./mocks/bridge.go -package=mocks_bindings
