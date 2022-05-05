package logger

import (
	"encoding/json"
)

type ExtLog struct {
	Level   string                 `json:"level"`
	Bridge  string                 `json:"bridge"`
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
	delete(extLog.Rest, "message")
	delete(extLog.Rest, "error")

	h.hook.Log(extLog)

	return len(p), nil
}
