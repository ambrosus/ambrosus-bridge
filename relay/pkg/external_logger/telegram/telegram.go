package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpClient interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type externalLogger struct {
	Token      string
	ChatId     int
	HttpClient HttpClient
	Prefix     string
}

func NewExternalLogger(token string, chatId int, prefix string, httpClient HttpClient) *externalLogger {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &externalLogger{
		Token:      token,
		ChatId:     chatId,
		HttpClient: httpClient,
		Prefix:     prefix,
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

func (t *externalLogger) LogError(msg string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Token)
	body := &request{
		ChatId:    t.ChatId,
		Text:      fmt.Sprintf("%s <b>We got an unexpected error:</b>\n%s", t.Prefix, msg),
		ParseMode: "html",
	}

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(body); err != nil {
		return fmt.Errorf("json encode request: %w", err)
	}
	resp, err := t.HttpClient.Post(url, "application/json", payloadBuf)
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
