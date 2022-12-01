package fee_helper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

type feeApiResponse struct {
	ThisGas string `json:"this"`
	SideGas string `json:"side"`
}

func getTransferFee(feeApiUrl string) (thisGas, sideGas decimal.Decimal, err error) {
	resp, err := http.Get(feeApiUrl)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, fmt.Errorf("failed to get transfer fee: %w", err)
	}
	defer resp.Body.Close()

	var r feeApiResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, fmt.Errorf("failed to parse transfer fee: %w", err)
	}

	thisGas, err = decimal.NewFromString(r.ThisGas)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, fmt.Errorf("failed to parse this gas: %w", err)
	}
	sideGas, err = decimal.NewFromString(r.SideGas)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, fmt.Errorf("failed to parse side gas: %w", err)
	}
	return thisGas, sideGas, nil
}
