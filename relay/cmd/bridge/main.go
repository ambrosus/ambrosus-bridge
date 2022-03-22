package main

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/eth"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger/telegram"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/metric"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Initialize bridge config.
	cfg, err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("error initialize config")
	}

	// Creating telegram loggers as an external logger.
	telegramAmbLogger := telegram.NewExternalLogger(cfg.Telegram.Token, cfg.Telegram.ChatId, "<b>[AMB]</b>", nil)
	telegramEthLogger := telegram.NewExternalLogger(cfg.Telegram.Token, cfg.Telegram.ChatId, "<b>[ETH]</b>", nil)

	// Creating a new ambrosus bridge.
	ambBridge, err := amb.New(&cfg.AMB, telegramAmbLogger)
	if err != nil {
		log.Fatal().Err(err).Msg("ambrosus bridge not created")
	}

	// Creating a new ethereum bridge.
	ethBridge, err := eth.New(&cfg.ETH, telegramEthLogger)
	if err != nil {
		log.Fatal().Err(err).Msg("ethereum bridge not created")
	}

	go ambBridge.Run(ethBridge)
	go ethBridge.Run(ambBridge)

	// Prometheus endpoint
	if err = metric.ServeEndpoint(cfg.Prometheus.Ip, cfg.Prometheus.Port); err != nil {
		log.Fatal().Err(err).Msg("failed to serve HTTP server (Prometheus endpoint)")
	}
}
