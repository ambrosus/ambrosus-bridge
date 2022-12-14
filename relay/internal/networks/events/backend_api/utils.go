package backend_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/events"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 50 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
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

	// ping
	pingDoneCh := make(chan struct{})
	pingErrCh := make(chan error)
	go ping(conn, pingErrCh, pingDoneCh)

	// resp
	type respErr struct {
		resp []byte
		err  error
	}
	respChan := make(chan respErr)
	go func() {
		_, resp, err := conn.ReadMessage()
		respChan <- respErr{resp, err}
	}()

	// listen the resp or errors from ping or ReadMessage
	select {
	case resp := <-respChan:
		close(pingDoneCh)
		if resp.err != nil {
			return nil, fmt.Errorf("read message: %w", resp.err)
		}
		return resp.resp, nil
	case err := <-pingErrCh:
		return nil, fmt.Errorf("write ping: %w", err)
	}
}

func ping(conn *websocket.Conn, errCh chan<- error, doneCh <-chan struct{}) {
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-doneCh:
			return
		case <-ticker.C:
			err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(pongWait))
			if err != nil {
				errCh <- err
				return
			}
		}
	}
}

func (a *EventsApi) waitEventUrl(eventName string) string {
	return fmt.Sprintf("ws://%s/ws/relay/waitEvent?networkFrom=%s&networkTo=%s&eventName=%s", a.baseUrl, a.thisName, a.sideName, eventName)
}

func (a *EventsApi) getEventUrl(eventName string, eventId uint64) string {
	return fmt.Sprintf("http://%s/relay/getEvent?networkFrom=%s&networkTo=%s&eventName=%s&eventId=%v", a.baseUrl, a.thisName, a.sideName, eventName, eventId)
}
