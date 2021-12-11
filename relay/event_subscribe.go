package main

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("")
	if err != nil {
		log.Fatal(err)
	}
	_ = client
}
