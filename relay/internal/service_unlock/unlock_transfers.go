package service_unlock

import (
	"context"
	"fmt"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_validity_watchdog"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

type UnlockTransfers struct {
	bridge        networks.Bridge
	watchValidity *service_validity_watchdog.WatchTransfersValidity
	logger        *zerolog.Logger
}

func NewUnlockTransfers(bridge networks.Bridge, watchValidity *service_validity_watchdog.WatchTransfersValidity) *UnlockTransfers {
	logger := bridge.GetLogger().With().Str("service", "UnlockTransfers").Logger()

	return &UnlockTransfers{
		bridge:        bridge,
		watchValidity: watchValidity,
		logger:        &logger,
	}
}

func (b *UnlockTransfers) Run() {
	cb.ShouldHavePk(b.bridge)
	for {
		cb.EnsureContractUnpaused(b.bridge, b.logger)

		if err := b.unlockOldTransfers(); err != nil {
			b.logger.Error().Err(err).Msg("")
		}
		time.Sleep(1 * time.Minute)
	}
}

func (b *UnlockTransfers) unlockOldTransfers() error {
	// Get oldest transfer timestamp.
	oldestLockedEventId, err := b.bridge.GetContract().OldestLockedEventId(nil)
	if err != nil {
		return fmt.Errorf("get oldest locked event id: %w", err)
	}
	lockedTransferTime, err := b.bridge.GetContract().LockedTransfers(nil, oldestLockedEventId)
	if err != nil {
		return fmt.Errorf("get locked transfer time %v: %w", oldestLockedEventId, err)
	}
	if lockedTransferTime.Uint64() == 0 {
		lockTime, err := b.bridge.GetContract().LockTime(nil)
		if err != nil {
			return fmt.Errorf("get lock time: %w", err)
		}

		b.logger.Debug().Str("event_id", oldestLockedEventId.String()).Msgf(
			"there are no locked transfers with that id. Sleep %v seconds...",
			lockTime.Uint64(),
		)
		time.Sleep(time.Duration(lockTime.Uint64()) * time.Second)
		return nil
	}

	// Get the latest block.
	latestBlock, err := b.bridge.GetClient().BlockByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("get latest block: %w", err)
	}

	// Check if the unlocking is allowed and get the sleep time.
	sleepTime := lockedTransferTime.Int64() - int64(latestBlock.Time())
	if sleepTime > 0 {
		b.logger.Debug().Str("event_id", oldestLockedEventId.String()).Msgf(
			"sleep %v seconds...",
			sleepTime,
		)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}

	// Unlock the oldest transfer.

	b.logger.Info().Str("event_id", oldestLockedEventId.String()).Msg("check validity of locked transfers...")
	if err := b.watchValidity.CheckOldLockedTransferFromId(oldestLockedEventId); err != nil {
		return fmt.Errorf("checkOldLockedTransferFromId: %w", err)
	}

	// b.logger.Info().Str("event_id", oldestLockedEventId.String()).Msgf("waiting next block before unlocking...")
	// if err := cb.WaitForNextBlock(b.bridge.GetWsClient()); err != nil {
	// 	return fmt.Errorf("wait for next block: %w", err)
	// }

	b.logger.Info().Str("event_id", oldestLockedEventId.String()).Msg("unlocking...")
	if err := b.unlockTransfers(); err != nil {
		return fmt.Errorf("unlock locked transfer %v: %w", oldestLockedEventId, err)
	}

	b.logger.Info().Str("event_id", oldestLockedEventId.String()).Msg("unlocked")
	return nil
}

func (b *UnlockTransfers) unlockTransfers() error {
	// Send the tx with the gas limit 60% more than the estimated gas limit.
	// todo: make the multiplier configurable
	authChangedGasLimit, err := cb.GetMultipliedEstimatedGasLimit(*b.bridge.GetAuth(), 1.60, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.bridge.GetContract().UnlockTransfersBatch(opts)
	})
	if err != nil {
		return fmt.Errorf("get multiplied estimated gas limit: %w", err)
	}

	return b.bridge.ProcessTx("unlockTransfersBatch", authChangedGasLimit, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.bridge.GetContract().UnlockTransfersBatch(opts)
	})
}
