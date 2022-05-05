package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type tgLogger struct {
	Token      string
	ChatId     int
	HttpClient *http.Client
	Prefix     string
}

func NewLogger(token string, chatId int, httpClient *http.Client) *tgLogger {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &tgLogger{
		Token:      token,
		ChatId:     chatId,
		HttpClient: httpClient,
	}
}

type request struct {
	ChatId    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type response struct {
	Ok               bool   `json:"ok"`
	ErrorDescription string `json:"description"` // if Ok is false
}

func (t *tgLogger) send(text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Token)
	body := &request{
		ChatId:    t.ChatId,
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

	respData := new(response)
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return fmt.Errorf("json decode response: %w", err)
	}
	if !respData.Ok {
		return fmt.Errorf(respData.ErrorDescription)
	}
	return nil

}
