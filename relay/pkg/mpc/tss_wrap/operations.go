package tss_wrap

import (
	"context"
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss-lib/common"
	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/tss"
)

// Keygen initiate share generation; Set received share to m.share after finishing protocol
func (m *Mpc) Keygen(
	ctx context.Context,
	inCh chan []byte,
	outCh chan *OutputMessage,
	errCh chan error,
	optionalPreParams ...keygen.LocalPreParams,
) {
	params := tss.NewParameters(tss.S256(), m.party, m.me, len(m.party.IDs()), len(m.party.IDs())-1)
	outChTss := make(chan tss.Message, 100)
	endCh := make(chan keygen.LocalPartySaveData, 5)
	keygenParty := keygen.NewLocalParty(params, outChTss, endCh, optionalPreParams...)

	go func() {
		if err := keygenParty.Start(); err != nil {
			errCh <- fmt.Errorf("keyGen.Start: %s", err.Error())
		}
	}()

	share, err := m.messaging(ctx, keygenParty, inCh, outCh, outChTss, endCh, nil)
	m.share = share.(*keygen.LocalPartySaveData)

	errCh <- err
}

// Sign initiate signing msg; Signature will be stored in `result` arg after finishing protocol
// nil is sent to encCh when signing is finished successfully
func (m *Mpc) Sign(
	ctx context.Context,
	inCh chan []byte,
	outCh chan *OutputMessage,
	errCh chan error,
	msg []byte,
	result *[]byte,
) {
	bigMsg := new(big.Int).SetBytes(msg)

	params := tss.NewParameters(tss.S256(), m.party, m.me, len(m.party.IDs()), len(m.party.IDs())-1)
	outChTss := make(chan tss.Message, 100)
	endCh := make(chan common.SignatureData, 5)

	signParty := signing.NewLocalParty(bigMsg, params, *m.share, outChTss, endCh)

	go func() {
		if err := signParty.Start(); err != nil {
			errCh <- fmt.Errorf("signParty.Start: %s", err.Error())
		}
	}()

	signature, err := m.messaging(ctx, signParty, inCh, outCh, outChTss, nil, endCh)
	errCh <- err
	*result = tssSignToECDSA(signature.(*common.SignatureData))
}

// messaging is a generic function for receiving and transmitting messages.
// loop ends when: ctx done OR some endCh receive result OR errCh receive error.
func (m *Mpc) messaging(
	ctx context.Context, party tss.Party,

	inCh chan []byte,
	outCh chan *OutputMessage,

	outChTss chan tss.Message,
	endChKeygen chan keygen.LocalPartySaveData,
	endChSign chan common.SignatureData,
) (interface{}, error) {
	for {
		select {
		case msgIn := <-inCh:
			parsedMessage, err := m.unmarshallInputMsg(msgIn)
			if err != nil {
				return nil, fmt.Errorf("unmarshallInputMsg: %s", err.Error())
			}
			ok, err := party.Update(parsedMessage)
			if !ok {
				return nil, fmt.Errorf("party.Update: %s", err.Error())
			}
			m.logger.Debug().Msgf("%v received msg from %v %v", m.me.Id, parsedMessage.GetFrom().Id, parsedMessage.Type())

		case msgOut := <-outChTss:
			outputMsg, err := m.newOutputMsg(msgOut)
			if err != nil {
				return nil, fmt.Errorf("newOutputMsg: %s", err.Error())
			}
			m.logger.Debug().Msgf("%v send messages to %v %v", m.me.Id, outputMsg.SendToIds, msgOut.Type())

			outCh <- outputMsg

		case result := <-endChKeygen:
			return &result, nil

		case result := <-endChSign:
			return &result, nil

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func tssSignToECDSA(signature *common.SignatureData) []byte {
	return append(signature.Signature, signature.SignatureRecovery[0])
}
