package tss_wrap

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// todo tests

type Mpc struct {
	me          *tss.PartyID
	party       *tss.PeerContext
	partyIDsMap map[string]*tss.PartyID
	threshold   int
	share       *keygen.LocalPartySaveData // share for local peer
}

type MpcConfig struct {
	meID      int
	partyLen  int
	threshold int
}

func NewMpc(cfg *MpcConfig) *Mpc {
	partyIDsMap, party := createParty(cfg.partyLen)
	return &Mpc{
		me:          partyIDsMap[fmt.Sprint(cfg.meID)],
		party:       party,
		partyIDsMap: partyIDsMap,
		threshold:   cfg.threshold,
	}
}

func (m *Mpc) Threshold() int {
	return m.threshold
}
func (m *Mpc) MyID() string {
	return m.me.Id
}

func (m *Mpc) Share() []byte {
	// todo export share for saving
	panic("not implemented")
}

// createParty party of partyLen participant; context == party
func createParty(partyLen int) (partyIDsMap map[string]*tss.PartyID, party *tss.PeerContext) {
	unsortedParty := make([]*tss.PartyID, partyLen)
	for id := 0; id < partyLen; id++ {
		stringID := fmt.Sprint(id)
		peer := tss.NewPartyID(stringID, stringID, big.NewInt(int64(id)))
		unsortedParty = append(unsortedParty, peer)
		partyIDsMap[stringID] = unsortedParty[id]
	}
	partyIDs := tss.SortPartyIDs(unsortedParty)
	return partyIDsMap, tss.NewPeerContext(partyIDs)
}

func (m *Mpc) GetAddress() (common.Address, error) {
	if m.share == nil {
		return common.Address{}, fmt.Errorf("peer has no share")
	}
	pubKey := ecdsa.PublicKey{
		Curve: secp256k1.S256(),
		X:     m.share.ECDSAPub.X(),
		Y:     m.share.ECDSAPub.Y(),
	}
	return crypto.PubkeyToAddress(pubKey), nil
}
