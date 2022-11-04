package networking

import (
	"context"
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

func TestFullMsg(t *testing.T) {
	server_ := server.NewServer(nil, "testtoken", &logger)
	ts := httptest.NewServer(server_)
	defer ts.Close()

	url := "ws://" + strings.TrimPrefix(ts.URL, "http://")
	client_ := client.NewClient(nil, url, "testtoken", &logger)

	server_.SetFullMsg([]byte("test"))
	msg, err := client_.GetFullMsg()
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), msg)
}

func TestNetworkingKeygen(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	partyIDs := []string{"0", "1", "2", "3", "4"}

	server_ := createServer(partyIDs[0], 5)
	ts := httptest.NewServer(server_)
	defer ts.Close()

	clients := createClients(partyIDs[1:], 5, ts.URL)

	waitForClients := doClientsOperation(clients, func(client_ *client.Client) {
		err := client_.Keygen(ctx, partyIDs)
		assert.NoError(t, err)
	})

	err := server_.Keygen(ctx, partyIDs)
	assert.NoError(t, err)
	waitForClients()

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

	clients := createClients(partyIDs[1:], 5, ts.URL)

	waitForClients := doClientsOperation(clients, func(client_ *client.Client) {
		_, err := client_.Sign(ctx, partyIDs, msg)
		assert.NoError(t, err)
	})

	signature, err := server_.Sign(ctx, partyIDs, msg)
	assert.NoError(t, err)
	waitForClients()
	// checks

	pubkey, err := server_.Tss.GetPublicKey()
	assert.NoError(t, err)

	sigPublicKey, err := crypto.SigToPub(msg, signature)
	assert.NoError(t, err)

	assert.Equal(t, pubkey, sigPublicKey)
}

func TestNetworkingRefresh(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	partyIDs := []string{"0", "1", "2", "3", "4", "backup"}

	server_ := createServer(partyIDs[0], 5)
	ts := httptest.NewServer(server_)
	defer ts.Close()

	clients := createClients(partyIDs[1:], 5, ts.URL)

	// keygen

	waitForClients := doClientsOperation(clients, func(client_ *client.Client) {
		err := client_.Keygen(ctx, partyIDs)
		assert.NoError(t, err)
	})

	err := server_.Keygen(ctx, partyIDs)
	assert.NoError(t, err)
	waitForClients()

	pubkeyOld, err := server_.Tss.GetPublicKey()
	assert.NoError(t, err)

	// refresh

	// peer 1 lost his share
	delete(clients, "1") // delete from networking
	oldPartyIDs := []string{"0", "2", "3", "4", "backup"}
	newPartyIDs := []string{"0a", "1a", "2a", "3a", "4a", "backup_a"}

	// add new party to peers map
	for id, peer := range createClients(newPartyIDs, 5, ts.URL) {
		clients[id] = peer
	}

	waitForClients = doClientsOperation(clients, func(client_ *client.Client) {
		err = client_.Reshare(ctx, oldPartyIDs, newPartyIDs, 5)
		assert.NoError(t, err)
	})

	err = server_.Reshare(ctx, oldPartyIDs, newPartyIDs, 5)
	assert.NoError(t, err)
	waitForClients()

	// checks

	pubkeyNew, err := server_.Tss.GetPublicKey()
	assert.NoError(t, err)
	assert.Equal(t, pubkeyOld, pubkeyNew)

}

func TestNetworkingAccessToken(t *testing.T) {
	server := server.NewServer(nil, "serverToken", nil)
	ts := httptest.NewServer(server)
	defer ts.Close()

	partyIDs := []string{"1"}
	clients := createClients(partyIDs, 1, ts.URL)

	err := clients["1"].Keygen(context.Background(), partyIDs)
	assert.Equal(t, "ws connect: websocket: bad handshake. http resp: 401 Unauthorized", err.Error())
}

func doClientsOperation(clients map[string]*client.Client, operation func(client_ *client.Client)) (waitFunc func()) {
	var wg sync.WaitGroup
	for _, client_ := range clients {
		wg.Add(1)
		go func(client_ *client.Client) {
			time.Sleep(time.Second) // wait for server to start operation
			operation(client_)
			wg.Done()
		}(client_)
	}
	return wg.Wait
}

func createServer(serverID string, threshold int) *server.Server {
	serverLogger := logger.With().Str("server", serverID).Logger()
	tssLogger := serverLogger.With().Str("tss", "").Logger()
	mpc, err := tss_wrap.NewMpcWithShare(serverID, threshold, fixtures.GetShare(serverID), &tssLogger)
	if err != nil {
		panic(err)
	}
	return server.NewServer(mpc, "testtoken", &serverLogger)
}

func createClients(clientsIDs []string, threshold int, httpUrl string) map[string]*client.Client {
	url := "ws://" + strings.TrimPrefix(httpUrl, "http://")
	clients := make(map[string]*client.Client)

	for _, id := range clientsIDs {
		clientLogger := logger.With().Str("client", id).Logger()
		tssLogger := clientLogger.With().Str("tss", "").Logger()
		mpc, err := tss_wrap.NewMpcWithShare(id, threshold, fixtures.GetShare(id), &tssLogger)
		if err != nil {
			panic(err)
		}

		client_ := client.NewClient(mpc, url, "testtoken", &clientLogger)
		clients[id] = client_
	}
	return clients
}
