package main

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram/middlewares"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/bsc"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/eth"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize bridge config.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error initialize config")
	}

	baseLogger := createLogger(cfg.ExtLoggers)

	ambBridge, sideBridge := createBridges(cfg.Networks, baseLogger)

	go runSubmitters(cfg.Submitters, ambBridge, sideBridge)
	go runWatchdogs(cfg.Watchdogs, ambBridge, sideBridge)
	go runUnlockers(cfg.Unlockers, ambBridge, sideBridge)
	go runTriggers(cfg.Triggers, ambBridge, sideBridge)
	go runFeeApi(cfg.FeeApi, ambBridge, sideBridge, &baseLogger)
	go runPrometheus(cfg.Prometheus)

	select {}

}

func createLogger(cfg *config.ExternalLoggers) zerolog.Logger {
	var tgLogger logger.Hook
	if tg := cfg.Telegram; tg.Enable {
		tgLogger = middlewares.NewAntiDoubleMiddleware(telegram.NewLogger(tg.Token, tg.ChatId, nil))
	}
	return logger.NewLoggerWithHook(tgLogger)
}

func createBridges(cfg *config.Networks, logger zerolog.Logger) (ambBridge *amb.Bridge, sideBridge service_submit.Receiver) {
	var err error

	// Creating a new ambrosus bridge.
	ambBridge, err = amb.New(cfg.Networks["amb"], logger)
	if err != nil {
		log.Fatal().Err(err).Msg("ambrosus bridge not created")
	}

	// Creating a side (eth or bsc) bridge.
	switch cfg.SideBridgeNetwork {
	case "ETH":
		sideBridge, err = eth.New(cfg.Networks["eth"], logger)
	case "BSC":
		sideBridge, err = bsc.New(cfg.Networks["bsc"], logger)
	default:
		log.Fatal().Msg("dunno which sideBridge to create")
	}
	if err != nil {
		log.Fatal().Err(err).Msgf("side (%v) bridge not created", cfg.SideBridgeNetwork)
	}

	return
}
