package fee_api

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const signatureFeeTimestamp = 30 * 60 // 30 minutes

type reqParams struct {
	TokenAddress common.Address `json:"tokenAddress"`
	IsAmb        bool           `json:"isAmb"`

	Amount           *hexutil.Big `json:"amount"`
	IsAmountWithFees bool         `json:"IsAmountWithFees"`
}

type result struct {
	BridgeFee   *hexutil.Big  `json:"bridgeFee"`
	TransferFee *hexutil.Big  `json:"transferFee"`
	Signature   hexutil.Bytes `json:"signature"`
}

func (p *FeeAPI) feesHandler(w http.ResponseWriter, r *http.Request) {
	var req reqParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.JSONError(w, NewAppError(nil, "error when decoding request body", err.Error()).Marshal(), http.StatusBadRequest)
		return
	}

	result, err := p.getFees(req)
	if err != nil {
		helpers.JSONError(w, err.Marshal(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (p *FeeAPI) getFees(req reqParams) (*result, *AppError) {
	var bridge, sideBridge networks.BridgeFeeApi
	var bridgePT BridgePriceTracker
	if req.IsAmb {
		bridge, sideBridge = p.amb, p.side
		bridgePT = p.ambPT
	} else {
		bridge, sideBridge = p.side, p.amb
		bridgePT = p.sidePT
	}

	// if token address is native, then it's "wrapWithdraw" and we need to work with "wrapperAddress"
	var tokenAddress = req.TokenAddress
	if req.TokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		var err error
		tokenAddress, err = p.getWrapperToken(bridge)
		if err != nil {
			return nil, NewAppError(nil, "error when getting wrapper address", err.Error())
		}
	}

	// get coin prices of this and side bridges for reusing it below
	thisCoinPrice, err := p.getNativeCoinPrice(bridge)
	if err != nil {
		return nil, NewAppError(nil, "error when getting this bridge coin price", err.Error())
	}
	sideCoinPrice, err := p.getNativeCoinPrice(sideBridge)
	if err != nil {
		return nil, NewAppError(nil, "error when getting side bridge coin price", err.Error())
	}

	// get the bridge fee (% from amount)
	bridgeFee, err := getBridgeFee(bridge, thisCoinPrice, p.cache, tokenAddress, (*big.Int)(req.Amount))
	if err != nil {
		return nil, NewAppError(nil, "error when getting bridge fee", err.Error())
	}

	// get the transfer fee
	transferFee, err := p.getTransferFee(bridge, bridgePT, thisCoinPrice, sideCoinPrice)
	if err != nil {
		return nil, NewAppError(nil, "error when getting transfer fee", err.Error())
	}

	amount := (*big.Int)(req.Amount)
	// if amount contains fees, then we need to remove them
	if req.IsAmountWithFees {
		amount.Sub(amount, new(big.Int).Add(bridgeFee, transferFee))
		if amount.Cmp(big.NewInt(0)) <= 0 {
			return nil, ErrAmountIsTooSmall
		}
	}

	// sign the price with private key
	message := buildMessage(tokenAddress, transferFee, bridgeFee, amount)
	signature, err := bridge.Sign(message)
	if err != nil {
		return nil, NewAppError(nil, "error when signing data", err.Error())
	}

	return &result{
		BridgeFee:   (*hexutil.Big)(bridgeFee),
		TransferFee: (*hexutil.Big)(transferFee),
		Signature:   signature,
	}, nil
}

func buildMessage(tokenAddress common.Address, transferFee, bridgeFee, amount *big.Int) []byte {
	// tokenAddress + timestamp + transferFee + bridgeFee

	var data bytes.Buffer
	var b32 [32]byte // solidity fills uint256 with zero

	data.Write(tokenAddress.Bytes())

	timestamp := time.Now().Unix() / signatureFeeTimestamp
	data.Write(big.NewInt(timestamp).FillBytes(b32[:]))

	data.Write(transferFee.FillBytes(b32[:]))
	data.Write(bridgeFee.FillBytes(b32[:]))
	data.Write(amount.FillBytes(b32[:]))

	return accounts.TextHash(crypto.Keccak256(data.Bytes()))
}

func (p *FeeAPI) getTransferFee(bridge networks.BridgeFeeApi, bridgePT BridgePriceTracker, thisCoinPrice, sideCoinPrice float64) (*big.Int, error) {
	// get gas cost per withdraw in side bridge currency
	gasCostInSideI, err, _ := p.cache.Memoize("gasPerWithdraw"+bridge.GetName(), func() (interface{}, error) {
		return bridgePT.GasPerWithdraw(), nil
	})
	if err != nil {
		return nil, err
	}

	gasCostInSide := gasCostInSideI.(*big.Int)
	if gasCostInSide == nil {
		return bridge.GetDefaultTransferFeeWei(), nil
	}

	// convert it to native bridge currency
	gasCostInUsd := Coin2Usd(gasCostInSide, sideCoinPrice, 18)
	gasCostInNative := Usd2Coin(gasCostInUsd, thisCoinPrice, 18)
	return gasCostInNative, nil
}

func (p *FeeAPI) getNativeCoinPrice(bridge networks.BridgeFeeApi) (float64, error) {
	coinPriceI, err, _ := p.cache.Memoize("coinPrice"+bridge.GetName(), func() (interface{}, error) {
		return bridge.CoinPrice()
	})
	if err != nil {
		return 0, err
	}
	return coinPriceI.(float64), nil
}

func (p *FeeAPI) getWrapperToken(bridge networks.BridgeFeeApi) (common.Address, error) {
	tokenAddressI, err, _ := p.cache.Memoize("wrapperAddress"+bridge.GetName(), func() (interface{}, error) {
		return bridge.GetWrapperAddress()
	})
	if err != nil {
		return common.Address{}, err
	}
	return tokenAddressI.(common.Address), nil
}
