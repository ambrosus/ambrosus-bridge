package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

func GetBalanceGWei(client *ethclient.Client, contractAddress common.Address) (float64, error) {
	balance, err := client.BalanceAt(context.Background(), contractAddress, nil)
	if err != nil {
		return 0, err
	}
	balanceGWei := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(params.GWei))
	balanceFloat64, _ := balanceGWei.Float64()
	return balanceFloat64, nil
}
