package telegram

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

var (
	client = &mockHttpClient{}
	logger = NewExternalLogger("", 0, "", client)
)

type mockHttpClient struct {
	PostFunc func(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

func (m *mockHttpClient) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	return m.PostFunc(url, bodyType, body)
}

func TestNewExternalLogger(t *testing.T) {
	type args struct {
		token      string
		chatId     int
		prefix     string
		httpClient HttpClient
	}
	tests := []struct {
		name string
		args args
		want *externalLogger
	}{
		{
			name: "HttpClient is default",
			args: args{httpClient: nil},
			want: &externalLogger{HttpClient: &http.Client{}},
		},
		{
			name: "HttpClient is custom",
			args: args{httpClient: client},
			want: &externalLogger{HttpClient: client},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExternalLogger(tt.args.token, tt.args.chatId, tt.args.prefix, tt.args.httpClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExternalLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestLogError(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr error
	}{
		{
			name:    "Success",
			body:    `{"ok": true}`,
			wantErr: nil,
		},
		{
			name:    "Error from telegram",
			body:    `{"ok": false, "description": "Bad Request: chat not found"}`,
			wantErr: fmt.Errorf("Bad Request: chat not found"),
		},
		{
			name:    "Error when decode response",
			body:    `wrong json`,
			wantErr: fmt.Errorf("json decode response: invalid character 'w' looking for beginning of value"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.PostFunc = func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(tt.body)),
				}, nil
			}
			if err := logger.LogError(`{`); (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("LogError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
