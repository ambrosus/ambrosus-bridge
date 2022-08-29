package ambrosus_explorer

import (
	"net/http"
	nurl "net/url"
)

const (
	MainnetAPIBase = "https://explorer-api.ambrosus.io"
	TestnetAPIBase = "https://explorer-api.ambrosus-test.io"
	DevnetAPIBase  = "https://explorer-api.ambrosus-dev.io"
)

type AmbrosusExplorer struct {
	apiUrl     nurl.URL
	httpClient *http.Client
}

func NewAmbrosusExplorer(url string, httpClient *http.Client) (*AmbrosusExplorer, error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	if httpClient == nil {
		httpClient = new(http.Client)
	}

	return &AmbrosusExplorer{
		apiUrl:     *u,
		httpClient: httpClient,
	}, nil
}
