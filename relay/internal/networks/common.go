package networks

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

type CommonBridge struct {
	Bridge
	Client      *ethclient.Client
	WsClient    *ethclient.Client
	Contract    *contracts.Bridge
	WsContract  *contracts.Bridge
	ContractRaw *contracts.BridgeRaw
	Auth        *bind.TransactOpts
	SideBridge  Bridge
	Logger      zerolog.Logger
}

