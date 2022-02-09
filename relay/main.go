package main

import (
	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/eth"
)

func main() {

	ambBridge := amb.New(config.Config["amb"])
	ethBridge := eth.New(config.Config["eth"])

	go ambBridge.Run(ethBridge, ethBridge.SubmitBlockPoA)
	go ethBridge.Run(ambBridge, ambBridge.SubmitBlockPoW)

}
