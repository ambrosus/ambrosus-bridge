package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
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
