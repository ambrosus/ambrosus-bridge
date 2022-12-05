package main

import (
	"github.com/ambrosus/ambrosus-bridge/relay/cmd"
	"github.com/ambrosus/ambrosus-bridge/relay/cmd/untrustlessMonitoring/service_monitoring"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/eth"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize bridge config.
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error initialize config")
	}
	monitoringCfgViper, err := config.LoadConfig("configs/monitoring.json")
	if err != nil {
		log.Fatal().Err(err).Msg("error loading monitoring config")
	}
	monitoringCfg := new(service_monitoring.MonitoringConfig)
	if err = monitoringCfgViper.Unmarshal(monitoringCfg); err != nil {
		log.Fatal().Err(err).Msg("error initialize monitoring config")
	}

	baseLogger := cmd.CreateLogger(cfg.ExtLoggers).With().Str("service", "SubmitterAura").Logger()

	ambBridge, sideBridge, err := cmd.CreateBridges(cfg.Networks, cfg.EventsApi, baseLogger)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating bridges")
	}
	if _, ok := sideBridge.(*eth.Bridge); !ok {
		log.Fatal().Msg("side bridge is not eth bridge")
	}

	monitoringService, err := service_monitoring.NewMonitoring(monitoringCfg, ambBridge, sideBridge, baseLogger)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating monitoring service")
	}
	monitoringService.Run()
}
