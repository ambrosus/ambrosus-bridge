package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type requestSend struct {
	ChatId    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type responseSend struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID uint64 `json:"message_id"`
	} `json:"result"`

	ErrorDescription string `json:"description"` // if Ok is false
}

// Send method will split too long message into parts (telegram limits that to 4096 symbols)
func (t *TgLogger) Send(text string) (ids []uint64, parts []string, err error) {
	partsText := safeSplitText(text, maxMessageLength, separator)
	for _, part := range partsText {
		id, err := t.send(part)
		if err != nil {
			return ids, partsText, err
		}
		ids = append(ids, id)
	}

	return ids, partsText, nil
}

func (t *TgLogger) send(text string) (uint64, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Token)
	body := &requestSend{
		ChatId:    t.ChatId,
		Text:      text,
		ParseMode: "html",
	}

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(body); err != nil {
		return 0, fmt.Errorf("json encode request: %w", err)
	}
	resp, err := http.Post(url, "application/json", payloadBuf)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respData := new(responseSend)
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return 0, fmt.Errorf("json decode response: %w", err)
	}
	if !respData.Ok {
		return 0, fmt.Errorf(respData.ErrorDescription)
	}
	return respData.Result.MessageID, nil

}
