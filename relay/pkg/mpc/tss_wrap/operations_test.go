package tss_wrap

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/fixtures"
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

	doOperation(peers,
		func(p *testPeer, outCh chan *Message) {
			err := p.keygen(outCh)
			if err != nil {
				t.Fatal(err)
			}
		},
	)

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

	msg := fixtures.Message()
	signatures := make(map[string][]byte)

	doOperation(peers,
		func(p *testPeer, outCh chan *Message) {
			signature, err := p.sign(outCh, msg)
			if err != nil {
				t.Fatal(p.peer.MyID(), err)
			}

			signatures[p.peer.MyID()] = signature
		},
	)

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

func doOperation(peers map[string]*testPeer, operation func(*testPeer, chan *Message)) {
	outCh := make(chan *Message, 1000)

	var wg sync.WaitGroup
	wg.Add(len(peers))

	for _, peer := range peers {
		go func(p *testPeer) {
			defer wg.Done()
			fmt.Println("starting operation", p.peer.MyID())
			operation(p, outCh)
			fmt.Println("operation done", p.peer.MyID())
		}(peer)
	}

	fmt.Println("starting messaging")
	go messaging(outCh, peers)

	wg.Wait()
	close(outCh)
}

func messaging(outCh chan *Message, peers map[string]*testPeer) {
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
			peer: NewMpc(i, count, &log.Logger),
			inCh: make(chan []byte, 100),
		}

		peers[peer.peer.MyID()] = peer
	}
	return peers
}

func (p *testPeer) sign(outCh chan *Message, msg []byte) ([]byte, error) {
	if err := p.peer.SetShare(fixtures.GetShare(p.peer.me.Index)); err != nil {
		return nil, err
	}
	return p.peer.Sign(context.Background(), p.inCh, outCh, msg)
}

func (p *testPeer) keygen(outCh chan *Message) error {
	return p.peer.Keygen(context.Background(), p.inCh, outCh, fixtures.GetPreParams(p.peer.me.Index))
}
