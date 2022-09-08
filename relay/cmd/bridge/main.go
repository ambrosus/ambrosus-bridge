package main

import (
	"github.com/ambrosus/ambrosus-bridge/relay/cmd"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize bridge config.
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error initialize config")
	}

	baseLogger := cmd.CreateLogger(cfg.ExtLoggers)

	ambBridge, sideBridge, err := cmd.CreateBridges(cfg.Networks, baseLogger)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating bridges")
	}

	go runSubmitters(cfg.Submitters, ambBridge, sideBridge, baseLogger)
	go runWatchdogs(cfg.Watchdogs, ambBridge, sideBridge, baseLogger)
	go runUnlockers(cfg.Unlockers, ambBridge, sideBridge, baseLogger)
	go runTriggers(cfg.Triggers, ambBridge, sideBridge, baseLogger)
	go runFeeApi(cfg.FeeApi, ambBridge, sideBridge, baseLogger)
	go runPrometheus(cfg.Prometheus, baseLogger)
	go runHealth(":80", baseLogger)

	select {}

}
