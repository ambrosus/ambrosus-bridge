package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	maxMessageLength = 4096
	separator        = " "
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

// Send method will split too long message into parts (telegram limits that to 4096 symbols)
func (t *tgLogger) Send(text string) error {
	partsText := safeSplitText(text, maxMessageLength, separator)
	for _, part := range partsText {
		if err := t.send(part); err != nil {
			return err
		}
	}
	return nil
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

// safeSplitText splits text into slices of maxLength by separator
func safeSplitText(text string, maxLength int, separator string) []string {
	tempText := text
	var slices []string

	for len(tempText) > 0 {

		if len(tempText) > maxLength {
			splitPosition := strings.LastIndex(tempText[:maxLength], separator)

			// if there is no separator in the last 1/4 (of maxLength) symbols then split at maxLength
			if splitPosition < maxLength*3/4 {
				splitPosition = maxLength
			}

			slices = append(slices, tempText[:splitPosition])

			tempText = tempText[splitPosition:]
			tempText = strings.TrimLeft(tempText, "\n\t\v\f\r ")

		} else {
			slices = append(slices, tempText)
			break
		}

	}
	return slices
}
