package etherscan

import (
	"fmt"
	"net/http"
	nurl "net/url"
	"time"

	"github.com/nanmu42/etherscan-api"
)

const DefaultHttpClientTimeout = 30 * time.Second

type Etherscan struct {
	client *etherscan.Client
}

func NewEtherscan(url string, httpClient *http.Client) (*Etherscan, error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	uWithoutParams := u.Scheme + "://" + u.Host + u.Path
	apiKey := u.Query().Get("apikey")
	if apiKey == "" {
		return nil, fmt.Errorf("param `apikey` not found in the url")
	}

	customized := etherscan.Customization{
		Client:  httpClient,
		Key:     apiKey,
		BaseURL: uWithoutParams + "?",
	}
	if customized.Client == nil {
		customized.Timeout = DefaultHttpClientTimeout
	}

	return &Etherscan{client: etherscan.NewCustomized(customized)}, nil
}
