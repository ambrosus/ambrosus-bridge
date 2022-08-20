package logger

import (
	"encoding/json"
	"strings"
)

type ExtLog struct {
	Level   string                 `json:"level"`
	Bridge  string                 `json:"bridge"`
	Service string                 `json:"service"`
	Message string                 `json:"message"`
	Error   string                 `json:"error"`
	Rest    map[string]interface{} `json:"-"`
}

type Hook interface {
	Log(*ExtLog)
}

type hook struct {
	hook Hook
}

func (h hook) Write(p []byte) (n int, err error) {
	extLog := new(ExtLog)
	if err := json.Unmarshal(p, extLog); err != nil {
		return 0, err
	}
	if err := json.Unmarshal(p, &extLog.Rest); err != nil {
		return 0, err
	}
	delete(extLog.Rest, "level")
	delete(extLog.Rest, "bridge")
	delete(extLog.Rest, "service")
	delete(extLog.Rest, "message")
	delete(extLog.Rest, "error")

	if checkIsShouldLog(extLog) {
		h.hook.Log(extLog)
	}

	return len(p), nil
}

func checkIsShouldLog(extLog *ExtLog) bool {
	if strings.Contains(extLog.Error, "websocket: close 1006 (abnormal closure): unexpected EOF") ||
		strings.Contains(extLog.Error, "connection reset by peer") ||
		strings.Contains(extLog.Error, "cannot retry err [http2: Transport received Server's graceful shutdown GOAWAY] after Request.Body was written; define Request.GetBody to avoid this error") ||
		strings.Contains(extLog.Error, "i/o timeout") {
		return false
	}

	return true
}
