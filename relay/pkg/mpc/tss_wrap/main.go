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

type Mpc struct {
	meID      string
	threshold int
	share     *keygen.LocalPartySaveData // share for local peer
	logger    *zerolog.Logger
}

func NewMpcWithShare(meID string, threshold int, share []byte, logger *zerolog.Logger) (*Mpc, error) {
	mpc := NewMpc(meID, threshold, logger)
	if share != nil {
		err := mpc.SetShare(share)
		return mpc, err
	}
	return mpc, nil
}

func NewMpc(meID string, threshold int, logger *zerolog.Logger) *Mpc {
	return &Mpc{
		meID:      meID,
		threshold: threshold,
		logger:    logger,
	}
}

func (m *Mpc) Threshold() int {
	return m.threshold
}
func (m *Mpc) MyID() string {
	return m.meID
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

func createPeer(id string) *tss.PartyID {
	return tss.NewPartyID(id, id, new(big.Int).SetBytes([]byte(id)))
}
