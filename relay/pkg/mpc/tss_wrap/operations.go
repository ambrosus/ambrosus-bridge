package tss_wrap

import (
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss-lib/common"
	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/tss"
)

// todo context

// Keygen initiate share generation; Set received share to m.share after finishing protocol
func (m *Mpc) Keygen(inCh chan []byte, outCh chan<- *OutputMessage) error {
	params := tss.NewParameters(tss.S256(), m.party, m.me, len(m.party.IDs()), m.threshold)
	outChTss := make(chan tss.Message)
	endCh := make(chan keygen.LocalPartySaveData)
	keygenParty := keygen.NewLocalParty(params, outChTss, endCh)

	err := keygenParty.Start()
	if err != nil {
		return fmt.Errorf("keyGen.Start: %s", err.Error())
	}
	for {
		select {
		case msgIn := <-inCh:
			parsedMessage, err := m.unmarshallInputMsg(msgIn)
			if err != nil {
				return fmt.Errorf("unmarshallInputMsg: %s", err.Error())
			}
			_, err = keygenParty.Update(parsedMessage)
			if err != nil {
				return fmt.Errorf("keygen.Update: %s", err.Error())
			}
		case msgOut := <-outChTss:
			outputMsg, err := m.newOutputMsg(msgOut)
			if err != nil {
				return fmt.Errorf("newOutputMsg: %s", err.Error())
			}
			outCh <- outputMsg
		case m.share = <-endCh:
			return nil
		}
	}
}

// Sign initiate signing msg; Return signature after finishing protocol
func (m *Mpc) Sign(inCh chan []byte, outCh chan<- *OutputMessage, msg []byte) ([]byte, error) {
	bigMsg := new(big.Int).SetBytes(msg)

	params := tss.NewParameters(tss.S256(), m.party, m.me, len(m.party.IDs()), m.threshold)
	outChTss := make(chan tss.Message)
	endCh := make(chan common.SignatureData)

	signParty := signing.NewLocalParty(bigMsg, params, *m.share, outChTss, endCh)

	err := signParty.Start()
	if err != nil {
		return nil, fmt.Errorf("signParty.Start: %s", err.Error())
	}

	for {
		select {
		case msgIn := <-inCh:
			parsedMessage, err := m.unmarshallInputMsg(msgIn)
			if err != nil {
				return nil, fmt.Errorf("unmarshallInputMsg: %s", err.Error())
			}
			_, err = signParty.Update(parsedMessage)
			if err != nil {
				return nil, fmt.Errorf("signParty.Update: %s", err.Error())
			}
		case msgOut := <-outChTss:
			outputMsg, err := m.newOutputMsg(msgOut)
			if err != nil {
				return nil, fmt.Errorf("newOutputMsg: %s", err.Error())
			}
			outCh <- outputMsg

		case signature := <-endCh:
			return tssSignToECDSA(&signature), nil
		}
	}
}

func tssSignToECDSA(signature *common.SignatureData) []byte {
	return append(signature.Signature, signature.SignatureRecovery[0])
}
