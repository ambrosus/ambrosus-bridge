package main

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/bsc"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/eth"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_api"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_helper"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/aura"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/posa"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/pow"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_trigger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_unlock"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_watchdog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func runSubmitters(cfg *config.Submitters, ambBridge *amb.Bridge, sideBridge service_submit.Receiver) {
	log.Info().Str("service", "submitter").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	auraSubmitter, err := aura.NewSubmitterAura(ambBridge, &aura.ReceiverAura{Receiver: sideBridge}, cfg.Aura)
	if err != nil {
		log.Fatal().Err(err).Msg("auraBridgeSubmitter don't created")
	}

	var sideBridgeSubmitter service_submit.Submitter
	switch sideBridge.(type) {
	case *eth.Bridge:
		sideBridgeSubmitter, err = pow.NewSubmitterPoW(sideBridge, &pow.ReceiverPoW{Receiver: ambBridge}, cfg.Pow)
	case *bsc.Bridge:
		sideBridgeSubmitter, err = posa.NewSubmitterPoSA(sideBridge, &posa.ReceiverPoSA{Receiver: ambBridge})
	}
	if err != nil {
		log.Fatal().Err(err).Msg("sideBridgeSubmitter don't created")
	}

	go service_submit.NewSubmitTransfers(auraSubmitter, sideBridge).Run()
	go service_submit.NewSubmitTransfers(sideBridgeSubmitter, ambBridge).Run()
}

func runWatchdogs(cfg *config.Watchdogs, ambBridge *amb.Bridge, sideBridge networks.Bridge) {
	log.Info().Str("service", "watchdog").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	go service_watchdog.NewWatchTransfers(ambBridge, sideBridge.GetContract()).Run()
	go service_watchdog.NewWatchTransfers(sideBridge, ambBridge.GetContract()).Run()
}

func runUnlockers(cfg *config.Unlockers, ambBridge *amb.Bridge, sideBridge networks.Bridge) {
	log.Info().Str("service", "watchdog").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	ambWatchdog := service_watchdog.NewWatchTransfers(ambBridge, sideBridge.GetContract())
	go service_unlock.NewUnlockTransfers(ambBridge, ambWatchdog).Run()
	sideWatchdog := service_watchdog.NewWatchTransfers(sideBridge, ambBridge.GetContract())
	go service_unlock.NewUnlockTransfers(ambBridge, sideWatchdog).Run()
}

func runTriggers(cfg *config.Triggers, ambBridge *amb.Bridge, sideBridge networks.Bridge) {
	log.Info().Str("service", "triggers").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	go service_trigger.NewTriggerTransfers(ambBridge).Run()
	go service_trigger.NewTriggerTransfers(sideBridge).Run()
}

func runFeeApi(cfg *config.FeeApi, ambBridge, sideBridge networks.Bridge, logger zerolog.Logger) {
	log.Info().Str("service", "fee api").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	feeAmb, err := fee_helper.NewFeeHelper(ambBridge, sideBridge, cfg.Amb)
	if err != nil {
		log.Fatal().Err(err).Msg("feeAmb not created")
	}
	feeSide, err := fee_helper.NewFeeHelper(sideBridge, ambBridge, cfg.Side)
	if err != nil {
		log.Fatal().Err(err).Msg("feeSide not created")
	}

	feeApi := fee_api.NewFeeAPI(feeAmb, feeSide, logger)
	feeApi.Run(cfg.Endpoint, cfg.Ip, cfg.Port)
}

func runPrometheus(cfg *config.Prometheus) {
	log.Info().Str("service", "prometheus").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	if err := metric.ServeEndpoint(cfg.Ip, cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("failed to serve HTTP server (Prometheus endpoint)")
	}
}
