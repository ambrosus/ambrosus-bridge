package main

import (
	"fmt"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/eth"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger/telegram"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/metric"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize bridge config.
	cfg, err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("error initialize config")
	}

	// Creating telegram loggers as an external logger.
	telegramAmbLogger := telegram.NewExternalLogger(&cfg.Telegram, "<b>[AMB]</b>", nil)
	telegramEthLogger := telegram.NewExternalLogger(&cfg.Telegram, "<b>[ETH]</b>", nil)

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
	addr := fmt.Sprintf("%s:%d", cfg.Prometheus.Ip, cfg.Prometheus.Port)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal().Err(err).Msg("failed to serve HTTP server (Prometheus endpoint)")
	}
}
