package main

import (
	"context"
	"fmt"
	"math/big"
	"relay/networks/amb"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

const url = "https://network.ambrosus.io"

func Test(t *testing.T) {
	client, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}

	block, err := client.BlockByNumber(context.Background(), big.NewInt(16021709))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(block.Hash())
}

func TestHeaderHash(t *testing.T) {
	number := big.NewInt(16021709)
	h, _ := amb.HeaderByNumber(number)

	t.Log(h.Hash(false).Hex())
	t.Log(h.Hash(true).Hex())
}
