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
	partyIDs := []string{"0", "1", "2", "3", "4"}
	peers := createPeers(partyIDs, 5)

	doOperation(peers,
		func(p *testPeer, outCh chan *Message) {
			err := p.keygen(outCh, partyIDs)
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
	partyIDs := []string{"0", "1", "2", "3", "4"}
	peers := createPeers(partyIDs, 5)

	msg := fixtures.Message()
	signatures := make(map[string][]byte)

	doOperation(peers,
		func(p *testPeer, outCh chan *Message) {
			signature, err := p.sign(outCh, partyIDs, msg)
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

func TestReshare(t *testing.T) {
	partyIDs := []string{"0", "1", "2", "3", "4", "backup"}
	peers := createPeers(partyIDs, 5) // threshold == peers - 1
	// keygen
	doOperation(peers,
		func(p *testPeer, outCh chan *Message) {
			err := p.keygen(outCh, partyIDs)
			if err != nil {
				t.Fatal(err)
			}
		},
	)

	oldAddress, err := peers["0"].peer.GetAddress()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("keygen finished")

	// reshare

	// peer 0 lost his share
	delete(peers, "0") // delete from networking
	oldPartyIDs := []string{"1", "2", "3", "4", "backup"}
	newPartyIDs := []string{"0a", "1a", "2a", "3a", "4a", "backup_a"}

	// add new party to peers map
	for id, peer := range createPeers(newPartyIDs, 5) {
		peers[id] = peer
	}

	doOperation(peers,
		func(p *testPeer, outCh chan *Message) {
			err := p.reshare(outCh, oldPartyIDs, newPartyIDs)
			if err != nil {
				t.Fatal(err)
			}
		},
	)

	// checks
	for _, peer := range peers {
		address, err := peer.peer.GetAddress()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, oldAddress, address)
	}

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

func createPeers(ids []string, threshold int) map[string]*testPeer {
	peers := make(map[string]*testPeer)

	for _, id := range ids {
		logger := log.With().Str("id", id).Logger()
		peer := &testPeer{
			peer: NewMpc(id, threshold, &logger),
			inCh: make(chan []byte, 100),
		}

		peers[peer.peer.MyID()] = peer
	}
	return peers
}

func (p *testPeer) sign(outCh chan *Message, party []string, msg []byte) ([]byte, error) {
	if err := p.peer.SetShare(fixtures.GetShare(p.peer.MyID())); err != nil {
		return nil, err
	}
	return p.peer.Sign(context.Background(), party, p.inCh, outCh, msg)
}

func (p *testPeer) keygen(outCh chan *Message, party []string) (err error) {
	preParams := fixtures.GetPreParams(p.peer.MyID())
	if preParams == nil {
		fmt.Printf("generating pre params for %s\n", p.peer.MyID())
		preParams, err = keygen.GeneratePreParams(5 * time.Minute)
		fmt.Printf("pre params for %s generated\n", p.peer.MyID())
		marshalled, err := json.Marshal(preParams)
		if err != nil {
			panic(err)
		}
		fmt.Print(string(marshalled))
	}
	return p.peer.Keygen(context.Background(), party, p.inCh, outCh, *preParams)
}

func (p *testPeer) reshare(outCh chan *Message, partyOld, partyNew []string) error {
	return p.peer.Reshare(context.Background(), partyOld, partyNew, p.inCh, outCh)
}
