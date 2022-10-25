package networking

import (
	"context"
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

func TestNetworkingKeygen(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// todo use pre params
	partyIDs := []string{"0", "1", "2", "3", "4"}

	server_ := createServer(partyIDs[0], 5)
	ts := httptest.NewServer(server_)
	defer ts.Close()
	wsURL := strings.TrimPrefix(ts.URL, "http://")

	clients := createClients(partyIDs[1:], 5, wsURL)

	for _, client_ := range clients {
		go func(client_ *client.Client) {
			time.Sleep(time.Second) // wait for server to start keygen operation
			err := client_.Keygen(ctx, partyIDs)
			if err != nil {
				t.Error(err)
			}
		}(client_)
	}

	err := server_.Keygen(ctx, partyIDs)
	assert.NoError(t, err)

	// checks

	pubkeyServer, err := server_.Tss.GetPublicKey()
	assert.NoError(t, err)
	pubkeyClient, err := clients["1"].Tss.GetPublicKey()
	assert.NoError(t, err)

	assert.Equal(t, pubkeyServer, pubkeyClient)
}

func TestNetworkingSigning(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msg := fixtures.Message()

	partyIDs := []string{"0", "1", "2", "3", "4"}

	server_ := createServer(partyIDs[0], 5)
	ts := httptest.NewServer(server_)
	defer ts.Close()
	wsURL := strings.TrimPrefix(ts.URL, "http://")

	clients := createClients(partyIDs[1:], 5, wsURL)

	for _, client_ := range clients {
		go func(client_ *client.Client) {
			time.Sleep(time.Second) // wait for server to start sign operation
			_, err := client_.Sign(ctx, partyIDs, msg)
			if err != nil {
				t.Error(err)
			}
		}(client_)
	}

	signature, err := server_.Sign(ctx, partyIDs, msg)
	assert.NoError(t, err)

	// checks

	pubkey, err := server_.Tss.GetPublicKey()
	assert.NoError(t, err)

	sigPublicKey, err := crypto.SigToPub(msg, signature)
	assert.NoError(t, err)

	assert.Equal(t, pubkey, sigPublicKey)
}

func TestNetworkingRefresh(t *testing.T) {
	// todo use pre params

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	partyIDs := []string{"0", "1", "2", "3", "4", "backup"}

	server_ := createServer(partyIDs[0], 5)
	ts := httptest.NewServer(server_)
	defer ts.Close()
	wsURL := strings.TrimPrefix(ts.URL, "http://")

	clients := createClients(partyIDs[1:], 5, wsURL)

	for _, client_ := range clients {
		go func(client_ *client.Client) {
			time.Sleep(time.Second) // wait for server to start keygen operation
			err := client_.Keygen(ctx, partyIDs)
			if err != nil {
				t.Error(err)
			}
		}(client_)
	}

	err := server_.Keygen(ctx, partyIDs)
	assert.NoError(t, err)

	pubkeyOld, err := server_.Tss.GetPublicKey()
	assert.NoError(t, err)

	// refresh

	// peer 1 lost his share
	delete(clients, "1") // delete from networking
	oldPartyIDs := []string{"0", "2", "3", "4", "backup"}
	newPartyIDs := []string{"0a", "1a", "2a", "3a", "4a", "backup_a"}

	// add new party to peers map
	for id, peer := range createClients(newPartyIDs, 5, wsURL) {
		clients[id] = peer
	}

	for _, client_ := range clients {
		go func(client_ *client.Client) {
			time.Sleep(time.Second) // wait for server to start keygen operation
			err := client_.Reshare(ctx, oldPartyIDs, newPartyIDs)
			if err != nil {
				t.Error(err)
			}
		}(client_)
	}

	err = server_.Reshare(ctx, oldPartyIDs, newPartyIDs)
	assert.NoError(t, err)

	// checks

	pubkeyNew, err := server_.Tss.GetPublicKey()
	assert.NoError(t, err)
	assert.Equal(t, pubkeyOld, pubkeyNew)

}

func createServer(serverID string, threshold int) *server.Server {
	serverLogger := logger.With().Str("server", serverID).Logger()
	tssLogger := serverLogger.With().Str("tss", "").Logger()
	mpc, err := tss_wrap.NewMpcWithShare(serverID, threshold, fixtures.GetShare(serverID), &tssLogger)
	if err != nil {
		panic(err)
	}
	return server.NewServer(mpc, &serverLogger)
}

func createClients(clientsIDs []string, threshold int, url string) map[string]*client.Client {
	clients := make(map[string]*client.Client)

	for _, id := range clientsIDs {
		clientLogger := logger.With().Str("client", id).Logger()
		tssLogger := clientLogger.With().Str("tss", "").Logger()
		mpc, err := tss_wrap.NewMpcWithShare(id, threshold, fixtures.GetShare(id), &tssLogger)
		if err != nil {
			panic(err)
		}

		client_ := client.NewClient(mpc, url, nil, &clientLogger)
		clients[id] = client_
	}
	return clients
}
