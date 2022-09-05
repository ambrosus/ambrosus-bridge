package service_monitoring

import (
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	loggermiddlewares "github.com/ambrosus/ambrosus-bridge/relay/internal/logger/middlewares"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram/antidouble"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
)

type Monitoring struct {
	cfg                      *MonitoringConfig
	bridge, sideBridge       networks.Bridge
	logger, monitoringLogger zerolog.Logger
	relayNames               map[common.Address]string
	minimumRequiredBalance   *big.Int
}

type MonitoringConfig struct {
	RelayNames                map[string]string       `json:"relayNames"`
	BalanceCheckTime          int                     `json:"balanceCheckTime"`
	ConfirmationCheckTime     int                     `json:"confirmationCheckTime"`
	ConfirmationTime          int                     `json:"confirmationTime"`
	MinimumRequiredBalanceWei string                  `json:"minimumRequiredBalanceWei"`
	LogTo                     *config.ExternalLoggers `mapstructure:"externalLogger"`
}

func NewMonitoring(cfg *MonitoringConfig, ambBridge, ethBridge networks.Bridge, logger zerolog.Logger) (*Monitoring, error) {
	relayNames := make(map[common.Address]string)
	for addr, name := range cfg.RelayNames {
		if !common.IsHexAddress(addr) {
			panic(fmt.Errorf("invalid address: %s (%v)", addr, name))
		}
		relayNames[common.HexToAddress(addr)] = name
	}
	monitoringLogger := createMonitoringLogger(cfg.LogTo)

	minimumRequiredBalance, ok := new(big.Int).SetString(cfg.MinimumRequiredBalanceWei, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert minimum required balance to big int (%s)", cfg.MinimumRequiredBalanceWei)
	}

	return &Monitoring{
		cfg:                    cfg,
		bridge:                 ambBridge,
		sideBridge:             ethBridge,
		logger:                 logger,
		monitoringLogger:       monitoringLogger,
		relayNames:             relayNames,
		minimumRequiredBalance: minimumRequiredBalance,
	}, nil
}

func (m *Monitoring) Run() {
	go m.balanceLoop()
	go m.confirmationsLoop()
	select {}
}

func createMonitoringLogger(cfg *config.ExternalLoggers) zerolog.Logger {
	var tgLoggerHook logger.Hook
	if tg := cfg.Telegram; tg.Enable {
		tgLogger := telegram.NewLogger(tg.Token, tg.ChatId, nil)
		tgLoggerHook = antidouble.NewAntiDoubleTgLogger(tgLogger)
		tgLoggerHook = loggermiddlewares.NewOnlyLevels(tgLoggerHook, []string{zerolog.LevelInfoValue})
	}
	return logger.NewLoggerWithHook(tgLoggerHook)
}

func (m *Monitoring) getRelays() ([]common.Address, error) {
	relays, err := m.bridge.GetContract().GetRelays(nil)
	if err != nil {
		return nil, fmt.Errorf("GetRelays: %w", err)
	}
	return relays, nil
}

func (m *Monitoring) formatAddress(relay common.Address) string {
	return fmt.Sprintf("%s (%s)", relay.Hex(), m.relayNames[relay])
}
