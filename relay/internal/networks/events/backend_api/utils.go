package backend_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/events"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gorilla/websocket"
)

var bridgeBoundContract = bind.NewBoundContract(common.Address{}, bindings.BridgeParsedABI, nil, nil, nil)

func parseToBinding(out interface{}, eventName string, data []byte) (types.Log, error) {
	log, err := parseLog(data)
	if err != nil {
		return log, fmt.Errorf("parse log: %w", err)
	}
	err = bridgeBoundContract.UnpackLog(out, eventName, log)
	if err != nil {
		return log, fmt.Errorf("unpack log: %w", err)
	}
	return log, nil
}

func parseLog(data []byte) (types.Log, error) {
	var log types.Log
	err := json.Unmarshal(data, &log)
	return log, err
}

func (a *EventsApi) get(eventName string, eventId uint64) ([]byte, error) {
	return get(a.getEventUrl(eventName, eventId))
}

func (a *EventsApi) wait(eventName string) ([]byte, error) {
	return wait(a.waitEventUrl(eventName))
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %v", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if len(data) == 0 {
		return nil, events.ErrEventNotFound
	}
	return data, nil
}

func wait(url string) ([]byte, error) {
	conn, httpResp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		if httpResp != nil {
			return nil, fmt.Errorf("ws connect: %w. http resp: %v", err, httpResp.Status)
		} else {
			return nil, fmt.Errorf("ws connect: %w", err)
		}
	}
	defer conn.Close()
	_, resp, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("read message: %w", err)
	}
	return resp, nil
}

func (a *EventsApi) waitEventUrl(eventName string) string {
	return fmt.Sprintf("ws://%s/relay/waitEvent?networkFrom=%s&networkTo=%s&eventName=%s", a.baseUrl, a.thisName, a.sideName, eventName)
}

func (a *EventsApi) getEventUrl(eventName string, eventId uint64) string {
	return fmt.Sprintf("http://%s/relay/getEvent?networkFrom=%s&networkTo=%s&eventName=%s&eventId=%v", a.baseUrl, a.thisName, a.sideName, eventName, eventId)
}
