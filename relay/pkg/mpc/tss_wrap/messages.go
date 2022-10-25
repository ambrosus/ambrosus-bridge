package tss_wrap

import (
	"encoding/json"
	"fmt"

	"github.com/bnb-chain/tss-lib/tss"
)

// Message should be sent to another participant(s)
type Message struct {
	SendToIds []string
	Message   []byte // marshaled inputMessage for another participant(s)
}

// inputMessage contains information that tss lib want to receive
type inputMessage struct {
	FromID       string
	FromIndex    int
	IsBroadcast  bool
	MsgWireBytes []byte
}

func (om *Message) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, om)
}

func (om *Message) Marshall() ([]byte, error) {
	return json.Marshal(om)
}

func (m *Mpc) newOutputMsg(msg tss.Message, allPeers []*tss.PartyID) (*Message, error) {
	// fetch IDs of all parties that should receive this message
	var sendToIds []string
	if msg.IsBroadcast() { // msg.GetTo will be nil for broadcast messages
		for _, peer := range allPeers {
			if peer.Id != m.meID {
				sendToIds = append(sendToIds, peer.Id)
			}
		}
	} else {
		for _, peer := range msg.GetTo() {
			if peer.Id == m.meID {
				panic("send to self ??")
			}
			sendToIds = append(sendToIds, peer.Id)
		}
	}

	// create and marshal message
	wireBytes, _, err := msg.WireBytes()
	if err != nil {
		return nil, fmt.Errorf("msg.WireBytes: %s", err.Error())
	}

	inputMsg := &inputMessage{
		FromID:       msg.GetFrom().GetId(),
		FromIndex:    msg.GetFrom().Index,
		IsBroadcast:  msg.IsBroadcast(),
		MsgWireBytes: wireBytes,
	}
	bytesToSend, err := json.Marshal(inputMsg)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %s", err.Error())
	}

	return &Message{
		SendToIds: sendToIds,
		Message:   bytesToSend,
	}, nil
}

func (m *Mpc) unmarshallInputMsg(msgWire []byte) (tss.ParsedMessage, error) {
	inputMsg := &inputMessage{}
	err := json.Unmarshal(msgWire, inputMsg)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	from := createPeer(inputMsg.FromID)
	from.Index = inputMsg.FromIndex
	message, err := tss.ParseWireMessage(inputMsg.MsgWireBytes, from, inputMsg.IsBroadcast)
	if err != nil {
		return nil, fmt.Errorf("tss.ParseWireMessage: %s", err.Error())
	}

	return message, nil
}
