package events

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

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
	return fmt.Sprintf("ws://%s/ws?networkFrom=%s&networkTo=%s&events=%s", a.baseUrl, a.thisName, a.sideName, eventName)
}

func (a *EventsApi) getEventUrl(eventName string, eventId uint64) string {
	return fmt.Sprintf("http://%s/relay/getEvent?networkFrom=%s&networkTo=%s&eventName=%s&eventId=%v", a.baseUrl, a.thisName, a.sideName, eventName, eventId)
}
