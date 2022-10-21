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
	inCh <-chan []byte,
	outCh chan<- *Message,
	errCh chan<- error,
	optionalPreParams ...keygen.LocalPreParams,
) {
	params := tss.NewParameters(tss.S256(), m.party, m.me, len(m.party.IDs()), len(m.party.IDs())-1)
	outChTss := make(chan tss.Message, 100)
	endCh := make(chan keygen.LocalPartySaveData, 5)
	keygenParty := keygen.NewLocalParty(params, outChTss, endCh, optionalPreParams...)

	share, err := m.messaging(ctx, keygenParty, inCh, outCh, outChTss, endCh, nil)
	if share != nil {
		m.share = share.(*keygen.LocalPartySaveData)
	}
	errCh <- err
}

// Sign initiate signing msg; Signature will be stored in `result` arg after finishing protocol
// nil is sent to encCh when signing is finished successfully
func (m *Mpc) Sign(
	ctx context.Context,
	inCh <-chan []byte,
	outCh chan<- *Message,
	errCh chan<- error,
	msg []byte,
	result *[]byte,
) {
	bigMsg := new(big.Int).SetBytes(msg)
	params := tss.NewParameters(tss.S256(), m.party, m.me, len(m.party.IDs()), len(m.party.IDs())-1)
	outChTss := make(chan tss.Message, 100)
	endCh := make(chan common.SignatureData, 5)
	signParty := signing.NewLocalParty(bigMsg, params, *m.share, outChTss, endCh)

	signature, err := m.messaging(ctx, signParty, inCh, outCh, outChTss, nil, endCh)
	if signature != nil {
		*result = tssSignToECDSA(signature.(*common.SignatureData))
	}
	errCh <- err

}

func (m *Mpc) KeygenSync(ctx context.Context, inCh <-chan []byte, outCh chan<- *Message, optionalPreParams ...keygen.LocalPreParams) error {
	errCh := make(chan error)
	go m.Keygen(ctx, inCh, outCh, errCh, optionalPreParams...)
	return <-errCh
}

func (m *Mpc) SignSync(ctx context.Context, inCh <-chan []byte, outCh chan<- *Message, msg []byte) ([]byte, error) {
	errCh := make(chan error)
	var result []byte
	go m.Sign(ctx, inCh, outCh, errCh, msg, &result)
	err := <-errCh
	return result, err
}

// messaging is a generic function for receiving and transmitting messages.
// loop ends when: ctx done OR some endCh receive result OR errCh receive error.
func (m *Mpc) messaging(
	ctx context.Context,
	party tss.Party,

	inCh <-chan []byte,
	outCh chan<- *Message,

	outChTss chan tss.Message,
	endChKeygen chan keygen.LocalPartySaveData,
	endChSign chan common.SignatureData,
) (interface{}, error) {
	// todo hard to read, maybe refactor

	partyErrCh := make(chan *tss.Error)
	go func() { partyErrCh <- party.Start() }()

	msgOutErrCh := make(chan error)

	// this goro will stop when outChTss is closed (below)
	// a separate goro is needed to not stop receiving messages from outChTss after receiving the result
	go func() {
		for msgOut := range outChTss {
			outputMsg, err := m.newOutputMsg(msgOut)
			if err != nil {
				msgOutErrCh <- fmt.Errorf("newOutputMsg: %s", err.Error())
			}
			m.logger.Debug().Msgf("%v send messages to %v %v", m.me.Id, outputMsg.SendToIds, msgOut.Type())

			outCh <- outputMsg
		}
		msgOutErrCh <- nil
	}()

	result, err := func() (interface{}, error) {
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

			case result := <-endChKeygen:
				return &result, nil
			case result := <-endChSign:
				return &result, nil

			case <-ctx.Done():
				return nil, ctx.Err()

			case err := <-msgOutErrCh:
				if err != nil {
					return nil, err
				}
			case err := <-partyErrCh:
				if err != nil {
					return nil, fmt.Errorf("party.Start: %w", err)
				}
			}
		}
	}()

	close(outChTss) // close outChTss to stop sending messages to outCh
	if err != nil {
		return nil, err
	}

	if err = <-msgOutErrCh; err != nil { // wait for all messages to be sent
		return nil, err
	}
	// it's important to wait for all messages to be sent before returning result
	// otherwise, outCh may be closed before all messages are sent
	return result, nil

}

func tssSignToECDSA(signature *common.SignatureData) []byte {
	return append(signature.Signature, signature.SignatureRecovery[0])
}
