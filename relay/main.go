package main

import (
	"relay/config"
	"relay/networks/amb"
	"relay/networks/eth"
)

func main() {

	ambBridge := amb.New(config.Config["amb"])
	ethBridge := eth.New(config.Config["eth"])

	go ambBridge.Run(ethBridge, ethBridge.SubmitBlockPoA)
	go ethBridge.Run(ambBridge, ambBridge.SubmitBlockPoW)

}
