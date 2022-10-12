package networking

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
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

func TestNetworkingKeygen(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// todo use pre params
	server_ := createServer(0)
	ts := httptest.NewServer(server_)
	defer ts.Close()
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	clients := createClients(0, 5, wsURL)

	var wg sync.WaitGroup
	wg.Add(4)
	for _, client_ := range clients {
		go func(client_ *client.Client) {
			defer wg.Done()
			time.Sleep(time.Second) // wait for server to start keygen operation
			err := client_.Keygen(ctx)
			if err != nil {
				t.Error(err)
				return
			}
		}(client_)
	}

	err := server_.Keygen(ctx)
	assert.NoError(t, err)

	wg.Wait() // wait for clients

	// checks

	pubkeyServer, err := server_.Tss.GetPublicKey()
	assert.NoError(t, err)
	pubkeyClient, err := clients[0].Tss.GetPublicKey()
	assert.NoError(t, err)

	assert.Equal(t, pubkeyServer, pubkeyClient)
}

func TestManyNetworkingSigning(t *testing.T) {
	for {
		TestNetworkingSigning(t)
		fmt.Println("===================================================================")
	}
}

func TestNetworkingSigning(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msg := fixtures.Message()

	server_ := createServer(0)
	ts := httptest.NewServer(server_)
	defer ts.Close()
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	clients := createClients(0, 5, wsURL)

	for _, client_ := range clients {
		go func(client_ *client.Client) {
			time.Sleep(time.Second) // wait for server to start sign operation
			_, err := client_.Sign(ctx, msg)
			if err != nil {
				t.Error(err)
				return
			}
		}(client_)
	}

	signature, err := server_.Sign(ctx, msg)
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
	mpc, err := tss_wrap.NewMpcWithShare(0, 5, fixtures.GetShare(0), &serverLogger)
	if err != nil {
		panic(err)
	}
	return server.NewServer(mpc, &serverLogger)
}

func createClients(serverID int, partyLen int, url string) []*client.Client {
	var clients []*client.Client

	for i := 0; i < partyLen; i++ {
		if i == serverID {
			continue
		}

		clientLogger := logger.With().Int("client", i).Logger()
		mpc, err := tss_wrap.NewMpcWithShare(i, 5, fixtures.GetShare(i), &clientLogger)
		if err != nil {
			panic(err)
		}

		client_ := client.NewClient(mpc, url, &clientLogger)
		clients = append(clients, client_)
	}
	return clients
}
