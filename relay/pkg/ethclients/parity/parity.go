package parity

import (
	"context"
	"math/big"

	common_ethclient "github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	common_ethclient.Client
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
	return &Client{Client: *common_ethclient.NewClient(c), c: c}
}

// Blockchain Access

// EstimateGas override coz openethereum doesn't give us a full error message at eth_estimateGas,
// so do eth_call method to get the full error message
func (ec *Client) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	gas, err := ec.Client.EstimateGas(ctx, call)
	if err == nil {
		return gas, nil
	}
	_, err = ec.Client.CallContract(ctx, call, nil)
	return 0, helpers.ParseError(err)
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (ec *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var result hexutil.Uint64
	err := ec.c.CallContext(ctx, &result, "parity_nextNonce", account)
	return uint64(result), err
}

// ParityHeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *Client) ParityHeaderByNumber(ctx context.Context, number *big.Int) (*Header, error) {
	var head *Header
	err := ec.c.CallContext(ctx, &head, "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return head, err
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}
