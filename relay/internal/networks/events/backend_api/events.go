package backend_api

import (
	"fmt"
	"strconv"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/rs/zerolog"
)

type EventsApi struct {
	baseUrl            string // without ws:// or https://
	thisName, sideName string
	logger             *zerolog.Logger
}

func NewEventsApi(baseUrl, thisName, sideName string, logger *zerolog.Logger) *EventsApi {
	return &EventsApi{
		baseUrl:  baseUrl,
		thisName: thisName,
		sideName: sideName,
		logger:   logger,
	}
}

func (a *EventsApi) WatchPaused() error {
	_, err := a.wait("Paused")
	return err
}

func (a *EventsApi) WatchUnpaused() error {
	_, err := a.wait("Unpaused")
	return err
}

func (a *EventsApi) GetTransfer(eventID uint64) (*bindings.BridgeTransfer, error) {
	resp, err := a.get("Transfer", eventID)
	if err != nil {
		return nil, fmt.Errorf("get transfer event: %w", err)
	}
	e := new(bindings.BridgeTransfer)
	if err := parseToBinding(e, "Transfer", resp); err != nil {
		return nil, err
	}
	return e, nil
}

func (a *EventsApi) WatchTransfer() (*bindings.BridgeTransfer, error) {
	resp, err := a.wait("Transfer")
	if err != nil {
		return nil, fmt.Errorf("wait transfer event: %w", err)
	}
	e := new(bindings.BridgeTransfer)
	if err := parseToBinding(e, "Transfer", resp); err != nil {
		return nil, err
	}
	return e, nil
}

func (a *EventsApi) WatchTransferSubmit() (*bindings.BridgeTransferSubmit, error) {
	resp, err := a.wait("TransferSubmit")
	if err != nil {
		return nil, fmt.Errorf("get transferSubmit event: %w", err)
	}
	e := new(bindings.BridgeTransferSubmit)
	if err := parseToBinding(e, "TransferSubmit", resp); err != nil {
		return nil, err
	}
	return e, nil
}

func (a *EventsApi) WatchTransferFinish() (*bindings.BridgeTransferFinish, error) {
	resp, err := a.wait("TransferFinish")
	if err != nil {
		return nil, fmt.Errorf("get transferFinish event: %w", err)
	}
	e := new(bindings.BridgeTransferFinish)
	if err := parseToBinding(e, "TransferFinish", resp); err != nil {
		return nil, err
	}
	return e, nil
}

// temporary, while can't use backend fee api
func (a *EventsApi) GetWithdrawsCount(fromEvent, toEvent uint64) (int, error) {
	url := fmt.Sprintf("http://%s/relay/withdrawsCount?networkFrom=%s&networkTo=%s", a.baseUrl, a.thisName, a.sideName)
	resp, err := get(url)
	if err != nil {
		return 0, fmt.Errorf("get withdraws count: %w", err)
	}
	withdrawCount, err := strconv.Atoi(string(resp))
	if err != nil {
		return 0, fmt.Errorf("parse withdraws count: %w", err)
	}
	return withdrawCount, nil
}
