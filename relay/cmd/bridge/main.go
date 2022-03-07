package main

import (
	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/eth"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger/telegram"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/metric"
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
	if err := metric.SetContractBalance(amb.ContractBalanceGWeiGauge, ambBridge.Client, ambBridge.Auth.From); err != nil {
		log.Fatal().Err(err).Msg("Failed to init metric for ambrosus bridge")
	}

	// Creating a new ethereum bridge.
	ethBridge, err := eth.New(&cfg.ETH, telegramEthLogger)
	if err != nil {
		log.Fatal().Err(err).Msg("ethereum bridge not created")
	}
	if err := metric.SetContractBalance(eth.ContractBalanceGWeiGauge, ethBridge.Client, ethBridge.Auth.From); err != nil {
		log.Fatal().Err(err).Msg("Failed to init metric for ambrosus bridge")
	}

	go ambBridge.Run(ethBridge)
	go ethBridge.Run(ambBridge)

	// не знаю как это правильно делается в го
	// но для временного решения пойдёть
	select {}
}
