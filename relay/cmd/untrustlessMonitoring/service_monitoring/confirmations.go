package service_monitoring

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ethereum/go-ethereum/common"
)

func (m *Monitoring) confirmationsLoop() {
	for {
		err := m.sendConfirmations()
		if err != nil {
			m.logger.Error().Err(err).Msg("error checking confirmations")
		}
		time.Sleep(time.Duration(m.cfg.ConfirmationCheckTime) * time.Second)
	}
}

func (m *Monitoring) sendConfirmations() error {
	confirmations, err := m.checkConfirmations()
	if err != nil {
		return fmt.Errorf("checkConfirmations: %w", err)
	}

	if confirmations == nil {
		// all events are confirmed
		return nil
	}

	text := "Not confirmed:\n"
	for addr, confirmed := range confirmations {
		if !confirmed {
			text += fmt.Sprintf(" - %s\n", m.formatAddress(addr))
		}
	}
	m.monitoringLogger.Info().Msg(text)
	return nil
}

func (m *Monitoring) checkConfirmations() (map[common.Address]bool, error) {
	// get event from side bridge that is not confirmed yet
	nonConfirmedEvent, err := m.getNotConfirmedEvent()
	if err != nil {
		return nil, fmt.Errorf("getNotConfirmedEvent: %w", err)
	} else if nonConfirmedEvent == nil {
		return nil, nil
	}

	// check if event is old enough
	isOldEnough, err := m.isEventOldEnough(nonConfirmedEvent)
	if err != nil {
		return nil, fmt.Errorf("isEventOldEnough: %w", err)
	} else if !isOldEnough {
		return nil, nil
	}

	// get relays that are not confirmed yet

	relays, err := m.getRelays()
	if err != nil {
		return nil, err
	}
	confirmations := make(map[common.Address]bool)

	for _, relay := range relays {
		confirmed, err := m.bridge.GetContract().IsConfirmedByRelay(nil, relay, nonConfirmedEvent.EventId, nonConfirmedEvent.Queue)
		if err != nil {
			return nil, fmt.Errorf("IsConfirmed(%v): %w", relay, err)
		}
		confirmations[relay] = confirmed
	}

	return confirmations, nil
}

func (m *Monitoring) getNotConfirmedEvent() (*bindings.BridgeTransfer, error) {
	lastEventId, err := m.bridge.GetContract().InputEventId(nil)
	if err != nil {
		return nil, fmt.Errorf("GetLastReceivedEventId: %w", err)
	}

	nextEventId := new(big.Int).Add(lastEventId, big.NewInt(1))
	nextEvent, err := cb.GetEventById(m.sideBridge, nextEventId)
	if errors.Is(err, networks.ErrEventNotFound) { // no more old events
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("getEventById on id %v: %w", nextEventId.String(), err)
	}
	return nextEvent, nil
}

func (m *Monitoring) isEventOldEnough(event *bindings.BridgeTransfer) (bool, error) {
	block, err := m.sideBridge.GetClient().BlockByHash(context.Background(), event.Raw.BlockHash)
	if err != nil {
		return false, fmt.Errorf("TransactionByHash(%v): %w", event.Raw.TxHash, err)
	}
	isOldEnough := time.Unix(int64(block.Time()), 0).
		Add(time.Duration(m.cfg.ConfirmationTime) * time.Second).
		Before(time.Now())
	return isOldEnough, nil
}
