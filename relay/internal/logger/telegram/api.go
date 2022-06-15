package telegram

import (
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
