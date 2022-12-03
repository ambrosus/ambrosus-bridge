package backend_api

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestGetTransferEvent(t *testing.T) {
	api := NewEventsApi("localhost:8080", "amb", "eth", &zerolog.Logger{})
	resp, err := api.GetTransfer(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
