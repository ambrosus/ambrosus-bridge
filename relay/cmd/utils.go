package cmd

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	loggermiddlewares "github.com/ambrosus/ambrosus-bridge/relay/internal/logger/middlewares"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram/antidouble"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/bsc"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/eth"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/rs/zerolog"
)

func CreateBridges(cfg *config.Networks, logger zerolog.Logger) (ambBridge *amb.Bridge, sideBridge service_submit.Receiver, err error) {
	// Creating a side (eth or bsc) bridge.
	switch cfg.SideBridgeNetwork {
	case "ETH":
		sideBridge, err = eth.New(cfg.Networks["eth"], amb.BridgeName, logger)
	case "BSC":
		sideBridge, err = bsc.New(cfg.Networks["bsc"], amb.BridgeName, logger)
	default:
		return nil, nil, fmt.Errorf("dunno which sideBridge to create")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("side (%v) bridge not created: %w", cfg.SideBridgeNetwork, err)
	}

	// Creating a new ambrosus bridge.
	ambBridge, err = amb.New(cfg.Networks["amb"], sideBridge.GetName(), logger)
	if err != nil {
		return nil, nil, fmt.Errorf("ambrosus bridge not created: %w", err)
	}

	return
}

func CreateLogger(cfg *config.ExternalLoggers) zerolog.Logger {
	var tgLoggerHook logger.Hook
	if tg := cfg.Telegram; tg.Enable {
		tgLogger := telegram.NewLogger(tg.Token, tg.ChatId, nil)
		tgLoggerHook = antidouble.NewAntiDoubleTgLogger(tgLogger)
		tgLoggerHook = loggermiddlewares.NewOnlyLevels(tgLoggerHook, loggermiddlewares.DefaultOnlyLevels)
		tgLoggerHook = loggermiddlewares.NewSkipSubstrings(tgLoggerHook, loggermiddlewares.DefaultSkipSubstrings)
		tgLoggerHook = loggermiddlewares.NewRemoveHTMLTags(tgLoggerHook, loggermiddlewares.DefaultApplicableRegexps)
	}
	return logger.NewLoggerWithHook(tgLoggerHook)
}
