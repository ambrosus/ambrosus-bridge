package networking

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/fixtures"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/client"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/server"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var logger = log.Logger

func TestNetworkingSigning(t *testing.T) {
	msg := fixtures.Message()

	server_ := createServer(0)
	ts := httptest.NewServer(server_)
	defer ts.Close()
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	clients := createClients(0, 5, wsURL)

	go server_.Run()

	for _, client_ := range clients {
		go func(client_ *client.Client) {
			time.Sleep(time.Second) // wait for server to start sign operation
			_, err := client_.Sign(msg)
			if err != nil {
				t.Error(err)
				return
			}
		}(client_)
	}

	signature, err := server_.Sign(msg)
	assert.NoError(t, err)

	// checks

	pubkey, err := server_.Tss.GetPublicKey()
	assert.NoError(t, err)

	sigPublicKey, err := crypto.SigToPub(msg, signature)
	assert.NoError(t, err)

	assert.Equal(t, pubkey, sigPublicKey)
}

func createServer(serverID int) *server.Server {
	serverLogger := logger.With().Int("server", serverID).Logger()
	server_ := server.NewServer(tss_wrap.NewMpc(&tss_wrap.MpcConfig{
		MeID:      0,
		PartyLen:  5,
		Threshold: 5,
	}, serverLogger), &serverLogger)
	err := server_.Tss.SetShare(fixtures.GetShare(0))
	if err != nil {
		panic(err)
	}
	return server_
}

func createClients(serverID int, partyLen int, url string) []*client.Client {
	var clients []*client.Client

	for i := 0; i < partyLen; i++ {
		if i == serverID {
			continue
		}

		clientLogger := logger.With().Int("client", i).Logger()
		client_ := client.NewClient(tss_wrap.NewMpc(&tss_wrap.MpcConfig{
			MeID:      i,
			PartyLen:  5,
			Threshold: 5,
		}, clientLogger), url, &clientLogger)
		err := client_.Tss.SetShare(fixtures.GetShare(i))
		if err != nil {
			panic(err)
		}

		clients = append(clients, client_)
	}
	return clients
}
