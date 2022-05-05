package main

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/bsc"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/eth"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger/telegram"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize bridge config.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error initialize config")
	}

	var tgAmbLogger, tgEthLogger external_logger.ExternalLogger

	if tg := cfg.ExtLoggers.Telegram; tg.Enable {
		// Creating telegram loggers as an external logger.
		tgAmbLogger = telegram.NewLogger(tg.Token, tg.ChatId, "<b>[AMB]</b>", nil)
		tgEthLogger = telegram.NewLogger(tg.Token, tg.ChatId, "<b>[SIDE?]</b>", nil) // todo
	}

	// Creating a new ambrosus bridge.
	ambBridge, err := amb.New(cfg.Networks.AMB, tgAmbLogger)
	if err != nil {
		log.Fatal().Err(err).Msg("ambrosus bridge not created")
	}

	// Creating a side (eth or bsc) bridge.
	var sideBridge networks.BridgeReceiveAura
	switch {
	case cfg.Networks.ETH != nil:
		sideBridge, err = eth.New(cfg.Networks.ETH, tgEthLogger)
		sideBridge.(*eth.Bridge).SetSideBridge(ambBridge)

	case cfg.Networks.BSC != nil:
		sideBridge, err = bsc.New(cfg.Networks.BSC, tgEthLogger)
		sideBridge.(*bsc.Bridge).SetSideBridge(ambBridge)

	}
	if err != nil {
		log.Fatal().Err(err).Msg("side bridge not created")
	}

	if cfg.IsRelay {
		go ambBridge.Run()
		go sideBridge.Run()
	}
	if cfg.IsWatchdog {
		go ambBridge.ValidityWatchdog()
		go sideBridge.ValidityWatchdog()
	}

	if cfg.Prometheus.Enable {
		// Prometheus endpoint
		if err = metric.ServeEndpoint(cfg.Prometheus.Ip, cfg.Prometheus.Port); err != nil {
			log.Fatal().Err(err).Msg("failed to serve HTTP server (Prometheus endpoint)")
		}
	} else {
		select {}
	}
}
