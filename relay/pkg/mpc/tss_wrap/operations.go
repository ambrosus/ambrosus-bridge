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
func (m *Mpc) Keygen(inCh chan []byte, outCh chan<- *OutputMessage, optionalPreParams ...keygen.LocalPreParams) error {
	params := tss.NewParameters(tss.S256(), m.party, m.me, len(m.party.IDs()), len(m.party.IDs())-1)
	outChTss := make(chan tss.Message, 100)
	endCh := make(chan keygen.LocalPartySaveData, 5)
	errorCh := make(chan error)
	keygenParty := keygen.NewLocalParty(params, outChTss, endCh, optionalPreParams...)

	go func() {
		if err := keygenParty.Start(); err != nil {
			errorCh <- fmt.Errorf("keyGen.Start: %s", err.Error())
		}
	}()
	go m.messaging(inCh, outCh, outChTss, errorCh, keygenParty)

	for {
		select {
		case share := <-endCh:
			m.share = &share
			return nil
		case err := <-errorCh:
			return fmt.Errorf("keygenParty: %s", err.Error())
		}
	}
}

// Sign initiate signing msg; Return signature after finishing protocol
func (m *Mpc) Sign(inCh chan []byte, outCh chan<- *OutputMessage, msg []byte) ([]byte, error) {
	bigMsg := new(big.Int).SetBytes(msg)

	params := tss.NewParameters(tss.S256(), m.party, m.me, len(m.party.IDs()), len(m.party.IDs())-1)
	outChTss := make(chan tss.Message, 100)
	endCh := make(chan common.SignatureData, 5)
	errorCh := make(chan error)

	signParty := signing.NewLocalParty(bigMsg, params, *m.share, outChTss, endCh)

	go func() {
		if err := signParty.Start(); err != nil {
			errorCh <- fmt.Errorf("signParty.Start: %s", err.Error())
		}
	}()
	go m.messaging(inCh, outCh, outChTss, errorCh, signParty)

	for {
		select {
		case signature := <-endCh:
			return tssSignToECDSA(&signature), nil
		case err := <-errorCh:
			return nil, fmt.Errorf("signParty: %s", err.Error())
		}
	}
}

func (m *Mpc) messaging(
	inCh chan []byte,
	outCh chan<- *OutputMessage,
	outChTss chan tss.Message,
	errCh chan error,
	party tss.Party,
) {
	for {
		select {
		case msgIn := <-inCh:
			parsedMessage, err := m.unmarshallInputMsg(msgIn)
			if err != nil {
				errCh <- fmt.Errorf("unmarshallInputMsg: %s", err.Error())
				return
			}
			ok, err := party.Update(parsedMessage)
			if !ok {
				errCh <- fmt.Errorf("party.Update: %s", err.Error())
				return
			}
			m.logger.Debug().Msgf("%v received msg from %v %v", m.me.Id, parsedMessage.GetFrom().Id, parsedMessage.Type())

		case msgOut := <-outChTss:
			outputMsg, err := m.newOutputMsg(msgOut)
			if err != nil {
				errCh <- fmt.Errorf("newOutputMsg: %s", err.Error())
				return
			}
			m.logger.Debug().Msgf("%v send messages to %v %v", m.me.Id, outputMsg.SendToIds, msgOut.Type())

			outCh <- outputMsg
		}
	}
}

func tssSignToECDSA(signature *common.SignatureData) []byte {
	return append(signature.Signature, signature.SignatureRecovery[0])
}
