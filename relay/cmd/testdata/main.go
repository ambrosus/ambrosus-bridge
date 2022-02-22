package main

import (
	"log"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/eth"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	bh := "0xe9bff194462f02de6ba677cba76657b8eed353cd0bd95785faace53f20fb705f"

	bridge, err := eth.New(&config.Bridge{Url: "https://mainnet.infura.io/v3/ab050ca98686478e9e9b06dfc3b2f069"})
	if err != nil {
		log.Fatal(err)
	}

	bridge.DisputeBlock(common.HexToHash(bh))
}
