package tss_wrap

import (
	"context"
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss-lib/common"
	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/ecdsa/resharing"
	"github.com/bnb-chain/tss-lib/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/tss"
)

// Keygen initiate share generation; Set received share to m.share after finishing protocol
func (m *Mpc) Keygen(ctx context.Context, party []string, inCh <-chan []byte, outCh chan<- *Message, optionalPreParams ...keygen.LocalPreParams) error {
	params, err := m.createParams(party)
	if err != nil {
		return err
	}
	if len(optionalPreParams) == 0 && m.share != nil { // can provide preParams through share
		optionalPreParams = append(optionalPreParams, m.share.LocalPreParams)
	}
	outChTss := make(chan tss.Message, 100)
	endCh := make(chan keygen.LocalPartySaveData, 5)
	keygenParty := keygen.NewLocalParty(params, outChTss, endCh, optionalPreParams...)

	share, err := m.messaging(ctx, keygenParty, inCh, outCh, outChTss, endCh, nil, params.Parties().IDs())
	if err != nil {
		return err
	}
	if share == nil {
		return fmt.Errorf("share is nil")
	}
	m.share = share.(*keygen.LocalPartySaveData)
	return nil
}

// Sign initiate signing msg; Signature will be stored in `result` arg after finishing protocol
// nil is sent to encCh when signing is finished successfully
func (m *Mpc) Sign(ctx context.Context, party []string, inCh <-chan []byte, outCh chan<- *Message, msg []byte) ([]byte, error) {
	bigMsg := new(big.Int).SetBytes(msg)
	params, err := m.createParams(party)
	if err != nil {
		return nil, err
	}
	if m.share == nil {
		return nil, fmt.Errorf("share is nil")
	}
	outChTss := make(chan tss.Message, 100)
	endCh := make(chan common.SignatureData, 5)
	signParty := signing.NewLocalParty(bigMsg, params, *m.share, outChTss, endCh)

	signature, err := m.messaging(ctx, signParty, inCh, outCh, outChTss, nil, endCh, params.Parties().IDs())
	if err != nil {
		return nil, err
	}
	if signature == nil {
		return nil, fmt.Errorf("signature is nil")
	}
	return tssSignToECDSA(signature.(*common.SignatureData)), nil
}

// Reshare initiate share regeneration; Set received share to m.share after finishing protocol
func (m *Mpc) Reshare(ctx context.Context, partyOld, partyNew []string, thresholdNew int, inCh <-chan []byte, outCh chan<- *Message, optionalPreParams ...keygen.LocalPreParams) error {
	params, err := m.createReshareParams(partyOld, partyNew, thresholdNew)
	if err != nil {
		return err
	}
	if m.share == nil {
		emptyShare := keygen.NewLocalPartySaveData(len(partyNew))
		m.share = &emptyShare
		if len(optionalPreParams) == 1 {
			m.share.LocalPreParams = optionalPreParams[0]
		}
	}
	outChTss := make(chan tss.Message, 100)
	endCh := make(chan keygen.LocalPartySaveData, 5)
	reshareParty := resharing.NewLocalParty(params, *m.share, outChTss, endCh)

	share, err := m.messaging(ctx, reshareParty, inCh, outCh, outChTss, endCh, nil, params.OldAndNewParties())
	if err != nil {
		return err
	}
	if share == nil {
		return fmt.Errorf("share is nil")
	}
	if !params.IsOldCommittee() { // if I'm in old committee, I don't need to save new share coz it useless for me
		m.share = share.(*keygen.LocalPartySaveData)
	}

	m.logger.Info().Msgf("reshare finished, new share: %v", m.share)

	return nil
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

	allPeers []*tss.PartyID,
) (interface{}, error) {
	// todo hard to read, maybe refactor

	partyErrCh := make(chan *tss.Error)
	go func() { partyErrCh <- party.Start() }()

	msgOutErrCh := make(chan error)

	// this goro will stop when outChTss is closed (below)
	// a separate goro is needed to not stop receiving messages from outChTss after receiving the result
	go func() {
		for msgOut := range outChTss {
			outputMsg, err := m.newOutputMsg(msgOut, allPeers)
			if err != nil {
				msgOutErrCh <- fmt.Errorf("newOutputMsg: %s", err.Error())
			}
			m.logger.Debug().Msgf("%v send messages to %v %v", m.MyID(), outputMsg.SendToIds, msgOut.Type())

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
				m.logger.Debug().Msgf("%v received msg from %v %v", m.MyID(), parsedMessage.GetFrom().Id, parsedMessage.Type())

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

func (m *Mpc) createParty(partyIDs []string) (*tss.PeerContext, *tss.PartyID) {
	unsortedParty := make([]*tss.PartyID, 0, len(partyIDs))
	var me *tss.PartyID

	for _, id := range partyIDs {
		peer := createPeer(id)
		unsortedParty = append(unsortedParty, peer)
		if id == m.MyID() {
			me = peer
		}
	}
	party := tss.SortPartyIDs(unsortedParty)
	return tss.NewPeerContext(party), me
}

func (m *Mpc) createParams(partyIDs []string) (*tss.Parameters, error) {
	party, me := m.createParty(partyIDs)
	if me == nil {
		return nil, fmt.Errorf("can't find myself in partyIDs")
	}
	params := tss.NewParameters(tss.S256(), party, me, len(party.IDs()), m.threshold-1)
	return params, nil
}

func (m *Mpc) createReshareParams(partyIDsOld, partyIDsNew []string, thresholdNew int) (*tss.ReSharingParameters, error) {
	var me *tss.PartyID

	partyOld, me := m.createParty(partyIDsOld)
	partyNew, meNew := m.createParty(partyIDsNew)
	if me == nil {
		me = meNew
		if me == nil {
			return nil, fmt.Errorf("can't find myself in partyIDs")
		}
	}

	params := tss.NewReSharingParameters(tss.S256(),
		partyOld, partyNew, me,
		len(partyOld.IDs()), m.threshold-1,
		len(partyNew.IDs()), thresholdNew-1)

	return params, nil
}
