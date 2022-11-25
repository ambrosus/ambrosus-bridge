package common

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	ethclient.Client
	c *rpc.Client
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*Client, error) {
	return DialContext(context.Background(), rawurl)
}

func DialContext(ctx context.Context, rawurl string) (*Client, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

// NewClient creates a client that uses the given RPC client.
func NewClient(c *rpc.Client) (client *Client) {
	return &Client{Client: *ethclient.NewClient(c), c: c}
}

type TxGasPriceResp struct {
	GasPrice *hexutil.Big `json:"gasPrice"`
}

func (ec *Client) TxGasPriceFromResponse(ctx context.Context, txHash common.Hash) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var json TxGasPriceResp
	err := ec.c.CallContext(ctx, &json, "eth_getTransactionByHash", txHash)
	return json.GasPrice.ToInt(), err
}
