package tss_wrap

import (
	"encoding/json"
	"fmt"

	"github.com/bnb-chain/tss-lib/tss"
)

// OutputMessage should be sent to another participant(s)
type OutputMessage struct {
	SendToIds []string
	Message   []byte // marshaled inputMessage for another participant(s)
}

// inputMessage contains information that tss lib want to receive
type inputMessage struct {
	FromID       string
	IsBroadcast  bool
	MsgWireBytes []byte
}

func (om *OutputMessage) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, om)
}

func (om *OutputMessage) Marshall() ([]byte, error) {
	return json.Marshal(om)
}

func (m *Mpc) newOutputMsg(msg tss.Message) (*OutputMessage, error) {
	// fetch IDs of all parties that should receive this message
	var sendToIds []string
	if msg.IsBroadcast() { // nsg.GetTo will be nil for broadcast messages
		for peerID := range m.partyIDsMap {
			if peerID != m.me.Id {
				sendToIds = append(sendToIds, peerID)
			}
		}
	} else {
		for _, peer := range msg.GetTo() {
			if peer.Id != m.me.Id {
				sendToIds = append(sendToIds, peer.Id)
			} else {
				panic("send to self ??")
			}
		}
	}

	// create and marshal message
	wireBytes, _, err := msg.WireBytes()
	if err != nil {
		return nil, fmt.Errorf("msg.WireBytes: %s", err.Error())
	}

	inputMsg := &inputMessage{
		FromID:       msg.GetFrom().GetId(),
		IsBroadcast:  msg.IsBroadcast(),
		MsgWireBytes: wireBytes,
	}
	bytesToSend, err := json.Marshal(inputMsg)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %s", err.Error())
	}

	return &OutputMessage{
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
	message, err := tss.ParseWireMessage(inputMsg.MsgWireBytes, m.partyIDsMap[inputMsg.FromID], inputMsg.IsBroadcast)
	if err != nil {
		return nil, fmt.Errorf("tss.ParseWireMessage: %s", err.Error())
	}

	return message, nil
}
