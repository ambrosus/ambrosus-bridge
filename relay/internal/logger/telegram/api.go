package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TgLogger struct {
	Token      string
	ChatId     int
	HttpClient *http.Client
	Prefix     string
}

func NewLogger(token string, chatId int, httpClient *http.Client) *TgLogger {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &TgLogger{
		Token:      token,
		ChatId:     chatId,
		HttpClient: httpClient,
	}
}

// send

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
	resp, err := t.api("sendMessage", &requestSend{
		ChatId:    t.ChatId,
		Text:      text,
		ParseMode: "html",
	})
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

// edit

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
	resp, err := t.api("editMessageText", &requestEditText{
		ChatId:    t.ChatId,
		MessageId: msgId,
		Text:      text,
		ParseMode: "html",
	})
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

// api

func (t *TgLogger) api(method string, params any) (*http.Response, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", t.Token, method)
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(params); err != nil {
		return nil, fmt.Errorf("json encode request: %w", err)
	}
	return t.HttpClient.Post(url, "application/json", payloadBuf)
}
