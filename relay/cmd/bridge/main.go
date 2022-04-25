package main

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/eth"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger/telegram"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/metric"
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
		tgEthLogger = telegram.NewLogger(tg.Token, tg.ChatId, "<b>[ETH]</b>", nil)
	}

	// Creating a new ambrosus bridge.
	ambBridge, err := amb.New(&cfg.Networks.AMB, tgAmbLogger)
	if err != nil {
		log.Fatal().Err(err).Msg("ambrosus bridge not created")
	}

	// Creating a new ethereum bridge.
	ethBridge, err := eth.New(&cfg.Networks.ETH, tgEthLogger)
	if err != nil {
		log.Fatal().Err(err).Msg("ethereum bridge not created")
	}

	if cfg.IsRelay {
		go ambBridge.Run(ethBridge)
		go ethBridge.Run(ambBridge)
	}
	if cfg.IsWatchdog {
		go ambBridge.ValidityWatchdog(ethBridge)
		go ethBridge.ValidityWatchdog(ambBridge)
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
