package tss_wrap

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/rs/zerolog"
)

// todo tests

type Mpc struct {
	me          *tss.PartyID
	party       *tss.PeerContext
	partyIDsMap map[string]*tss.PartyID
	threshold   int
	share       *keygen.LocalPartySaveData // share for local peer
	logger      zerolog.Logger
}

type MpcConfig struct {
	MeID      int
	PartyLen  int
	Threshold int
}

func NewMpc(cfg *MpcConfig, logger zerolog.Logger) *Mpc {
	partyIDsMap, party := createParty(cfg.PartyLen)
	return &Mpc{
		me:          partyIDsMap[fmt.Sprint(cfg.MeID)],
		party:       party,
		partyIDsMap: partyIDsMap,
		logger:      logger,
	}
}

func (m *Mpc) Threshold() int {
	return m.threshold
}
func (m *Mpc) MyID() string {
	return m.me.Id
}

func (m *Mpc) Share() ([]byte, error) {
	return json.Marshal(m.share)
}

func (m *Mpc) SetShare(share []byte) error {
	err := json.Unmarshal(share, &m.share)
	if err != nil {
		return err
	}
	if !m.share.ValidateWithProof() {
		return fmt.Errorf("invalid share")
	}
	return nil
}

// createParty party of partyLen participant; context == party
func createParty(partyLen int) (partyIDsMap map[string]*tss.PartyID, party *tss.PeerContext) {
	unsortedParty := make([]*tss.PartyID, 0, partyLen)
	partyIDsMap = make(map[string]*tss.PartyID)

	for id := 0; id < partyLen; id++ {
		stringID := fmt.Sprint(id)
		peer := tss.NewPartyID(stringID, stringID, big.NewInt(int64(id+1)))
		unsortedParty = append(unsortedParty, peer)
		partyIDsMap[stringID] = unsortedParty[id]
	}
	partyIDs := tss.SortPartyIDs(unsortedParty)
	return partyIDsMap, tss.NewPeerContext(partyIDs)
}

func (m *Mpc) GetAddress() (common.Address, error) {
	pubKey, err := m.GetPublicKey()
	return crypto.PubkeyToAddress(*pubKey), err
}

func (m *Mpc) GetPublicKey() (*ecdsa.PublicKey, error) {
	if m.share == nil {
		return nil, fmt.Errorf("peer has no share")
	}
	return &ecdsa.PublicKey{
		Curve: secp256k1.S256(),
		X:     m.share.ECDSAPub.X(),
		Y:     m.share.ECDSAPub.Y(),
	}, nil
}
