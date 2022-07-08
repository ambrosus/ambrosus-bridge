package main

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram/middlewares"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/bsc"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/eth"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_api"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_helper"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize bridge config.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error initialize config")
	}

	var tgLogger logger.Hook
	if tg := cfg.ExtLoggers.Telegram; tg.Enable {
		tgLogger = middlewares.NewAntiDoubleMiddleware(telegram.NewLogger(tg.Token, tg.ChatId, nil))
	}
	baseLogger := logger.NewLoggerWithHook(tgLogger)

	// Creating a new ambrosus bridge.
	ambBridge, err := amb.New(cfg.Networks.AMB, baseLogger)
	if err != nil {
		log.Fatal().Err(err).Msg("ambrosus bridge not created")
	}

	// Creating a side (eth or bsc) bridge.
	var sideBridge networks.BridgeReceiveAura
	switch {
	case cfg.Networks.ETH != nil:
		sideBridge, err = eth.New(cfg.Networks.ETH, baseLogger)
		if err != nil {
			log.Fatal().Err(err).Msg("eth bridge not created")
		}
		sideBridge.(*eth.Bridge).SetSideBridge(ambBridge)

	case cfg.Networks.BSC != nil:
		sideBridge, err = bsc.New(cfg.Networks.BSC, baseLogger)
		if err != nil {
			log.Fatal().Err(err).Msg("bsc bridge not created")
		}
		sideBridge.(*bsc.Bridge).SetSideBridge(ambBridge)
	}

	ambBridge.SetSideBridge(sideBridge)

	if cfg.IsRelay {
		go ambBridge.Run()
		go sideBridge.Run()
	}
	if cfg.IsWatchdog {
		go ambBridge.ValidityWatchdog()
		go sideBridge.ValidityWatchdog()
	}

	if feeCfg := cfg.FeeApi; feeCfg.Enable {
		feeAmb, err := fee_helper.NewFeeHelper(ambBridge, sideBridge, feeCfg.Amb)
		if err != nil {
			log.Fatal().Err(err).Msg("feeAmb not created")
		}
		feeSide, err := fee_helper.NewFeeHelper(sideBridge, ambBridge, feeCfg.Side)
		if err != nil {
			log.Fatal().Err(err).Msg("feeSide not created")
		}

		feeApi := fee_api.NewFeeAPI(feeAmb, feeSide, baseLogger)
		go feeApi.Run(cfg.FeeApi.Endpoint, cfg.FeeApi.Ip, cfg.FeeApi.Port)
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
