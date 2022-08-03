package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type requestEditText struct {
	ChatId    int    `json:"chat_id"`
	MessageId uint64 `json:"message_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type responseEditText struct {
	Ok               bool   `json:"ok"`
	ErrorDescription string `json:"description"` // if Ok is false
}

func (t *TgLogger) EditText(msgId uint64, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText", t.Token)
	body := &requestEditText{
		ChatId:    t.ChatId,
		MessageId: msgId,
		Text:      text,
		ParseMode: "html",
	}

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(body); err != nil {
		return fmt.Errorf("json encode request: %w", err)
	}
	resp, err := http.Post(url, "application/json", payloadBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respData := new(responseEditText)
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return fmt.Errorf("json decode response: %w", err)
	}
	if !respData.Ok {
		return fmt.Errorf(respData.ErrorDescription)
	}
	return nil

}
