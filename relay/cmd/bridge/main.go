package main

import (
	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/eth"
)

func main() {
	// Creating a new ambrosus bridge.
	ambBridge := amb.New(config.Config["amb"])
	// Creating a new ethereum bridge.
	ethBridge := eth.New(config.Config["eth"])

	go ambBridge.Run(ethBridge, ethBridge.SubmitBlockPoA)
	go ethBridge.Run(ambBridge, ambBridge.SubmitBlockPoW)
}
