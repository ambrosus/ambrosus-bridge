package tss_wrap

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/fixtures"
	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type testPeer struct {
	peer *Mpc
	inCh chan []byte
}

func TestKeygen(t *testing.T) {
	peers := createPeers(5)

	outCh := make(chan *OutputMessage, 1000)

	var wg sync.WaitGroup
	wg.Add(5)

	for _, peer := range peers {

		go func(p *testPeer) {
			defer wg.Done()
			fmt.Println("starting keygen", p.peer.MyID())

			err := p.keygen(outCh)
			fmt.Println("keygen done", p.peer.MyID())
			if err != nil {
				t.Error(err)
			}
		}(peer)
	}

	// messaging

	fmt.Println("starting messaging")
	go messaging(outCh, peers)

	wg.Wait()
	close(outCh)

	// checks

	address0, err := peers["0"].peer.GetAddress()
	if err != nil {
		t.Fatal(err)
	}

	for _, peer := range peers {
		address, err := peer.peer.GetAddress()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, address0, address)

		share, err := peer.peer.Share()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%v %v \n", peer.peer.MyID(), string(share))
	}
}

func TestSign(t *testing.T) {
	peers := createPeers(5)

	outCh := make(chan *OutputMessage, 1000)
	msg := fixtures.Message()

	signatures := make(map[string][]byte)

	var wg sync.WaitGroup
	wg.Add(5)

	for _, peer := range peers {

		go func(p *testPeer) {
			defer wg.Done()
			fmt.Println("starting signing", p.peer.MyID())

			signature, err := p.sign(outCh, msg)
			fmt.Println("signing done", p.peer.MyID())
			if err != nil {
				t.Error(p.peer.MyID(), err)
			}

			signatures[p.peer.MyID()] = signature
		}(peer)
	}

	// messaging

	fmt.Println("starting messaging")
	go messaging(outCh, peers)

	wg.Wait()
	close(outCh)

	// checks

	for _, signature := range signatures {
		assert.Equal(t, signatures["0"], signature)
	}

	signature := signatures["0"]
	pubkey, err := peers["0"].peer.GetPublicKey()
	assert.NoError(t, err)

	sigPublicKey, err := crypto.SigToPub(msg, signature)
	assert.NoError(t, err)

	assert.Equal(t, pubkey, sigPublicKey)
}

func TestGeneratePreParams(t *testing.T) {
	pp, err := keygen.GeneratePreParams(5 * time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	marshalled, err := json.Marshal(pp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(string(marshalled))
}

func messaging(outCh chan *OutputMessage, peers map[string]*testPeer) {
	for msg := range outCh {
		for _, id := range msg.SendToIds {
			//fmt.Println("send to", id)
			peers[id].inCh <- msg.Message
		}
	}
}

func createPeers(count int) map[string]*testPeer {
	peers := make(map[string]*testPeer)

	for i := 0; i < count; i++ {
		peer := &testPeer{
			peer: NewMpc(&MpcConfig{
				MeID:      i,
				PartyLen:  count,
				Threshold: count,
			}, log.Logger),
			inCh: make(chan []byte, 100),
		}

		peers[peer.peer.MyID()] = peer
	}
	return peers
}

func (p *testPeer) sign(outCh chan *OutputMessage, msg []byte) ([]byte, error) {
	if err := p.peer.SetShare(fixtures.GetShare(p.peer.me.Index)); err != nil {
		return nil, err
	}

	var signature []byte
	errCh := make(chan error, 1000)

	go p.peer.Sign(context.Background(), p.inCh, outCh, errCh, msg, &signature)

	err := <-errCh
	return signature, err
}

func (p *testPeer) keygen(outCh chan *OutputMessage) error {
	errCh := make(chan error, 1000)

	go p.peer.Keygen(context.Background(), p.inCh, outCh, errCh, fixtures.GetPreParams(p.peer.me.Index))

	err := <-errCh
	return err
}
