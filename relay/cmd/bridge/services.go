package main

import (
	"net/http"
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/api"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_helper"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_pause_unpause_watchdog"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/aura"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/posa"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/untrustless"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/untrustless2"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/untrustless_mpc"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_trigger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_unlock"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_validity_watchdog"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/client"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/server"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/rs/zerolog"
)

func runSubmitters(cfg *config.Submitters, ambBridge *amb.Bridge, sideBridge service_submit.Receiver, logger zerolog.Logger) {
	logger.Info().Str("service", "submitter").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	var err error

	// create amb->side submitter
	var ambSubmitter service_submit.Submitter
	switch cfg.AmbToSide {
	case "aura":
		ambSubmitter, err = aura.NewSubmitterAura(ambBridge, &aura.ReceiverAura{Receiver: sideBridge}, cfg.Aura)
	case "untrustless2":
		ambSubmitter, err = untrustless2.NewSubmitterUntrustless(ambBridge, &untrustless2.ReceiverUntrustless2{Receiver: sideBridge})
	case "untrustless-mpc":
		share, err := os.ReadFile(os.Getenv("SHARE_PATH"))
		if err != nil {
			logger.Fatal().Err(err).Msg("can't read share")
		}
		mpcc, err := tss_wrap.NewMpcWithShare(cfg.Mpc.MeID, cfg.Mpc.PartyLen, share, &logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to create MPC client")
		}

		var mpcSigner untrustless_mpc.MpcSigner
		if cfg.Mpc.IsServer {
			server_ := server.NewServerWithFullMsgEndpoint(mpcc, "/tx", &logger) // TODO: get full msg endpoint from config?
			go http.ListenAndServe(cfg.Mpc.ServerURL, server_)
			mpcSigner = server_
		} else {
			client_ := client.NewClient(mpcc, cfg.Mpc.ServerURL, &logger)
			mpcSigner = client_
		}

		receiver, err := untrustless_mpc.NewReceiverUntrustlessMpc(sideBridge, mpcSigner, "http://localhost:8082/tx", nil) // TODO: get serverTxUrl from config?
		ambSubmitter, err = untrustless_mpc.NewSubmitterUntrustlessMpc(ambBridge, receiver, cfg.Mpc.IsServer)

	default:
		logger.Info().Msg("amb->side submitter is disabled")
	}
	if err != nil {
		logger.Fatal().Err(err).Msg("ambBridgeSubmitter don't created")
	}

	// create side->amb submitter
	var sideBridgeSubmitter service_submit.Submitter
	switch cfg.SideToAmb {
	case "posa":
		sideBridgeSubmitter, err = posa.NewSubmitterPoSA(sideBridge, &posa.ReceiverPoSA{Receiver: ambBridge}, cfg.Posa)
	case "untrustless":
		sideBridgeSubmitter, err = untrustless.NewSubmitterUntrustless(sideBridge, &untrustless.ReceiverUntrustless{Receiver: ambBridge})
	default:
		logger.Info().Msg("side->amb submitter is disabled")
	}
	if err != nil {
		logger.Fatal().Err(err).Msg("sideBridgeSubmitter don't created")
	}

	// run submitters
	if ambSubmitter != nil {
		go service_submit.NewSubmitTransfers(ambSubmitter).Run()
	}
	if sideBridgeSubmitter != nil {
		go service_submit.NewSubmitTransfers(sideBridgeSubmitter).Run()
	}

}

func runValidityWatchdogs(cfg *config.ValidityWatchdogs, ambBridge *amb.Bridge, sideBridge networks.Bridge, logger zerolog.Logger) {
	logger.Info().Str("service", "validity watchdogs").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	if cfg.EnableForAmb {
		go service_validity_watchdog.NewWatchTransfersValidity(ambBridge, sideBridge.GetContract()).Run()
	}

	if cfg.EnableForSide {
		go service_validity_watchdog.NewWatchTransfersValidity(sideBridge, ambBridge.GetContract()).Run()
	}

}

func runPauseUnpauseWatchdogs(cfg *config.PauseUnpauseWatchdogs, ambBridge *amb.Bridge, sideBridge networks.Bridge, logger zerolog.Logger) {
	logger.Info().Str("service", "pause unpause watchdogs").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	go service_pause_unpause_watchdog.NewWatchPauseUnpause(ambBridge).Run()
	go service_pause_unpause_watchdog.NewWatchPauseUnpause(sideBridge).Run()

}

func runUnlockers(cfg *config.Unlockers, ambBridge *amb.Bridge, sideBridge networks.Bridge, logger zerolog.Logger) {
	logger.Info().Str("service", "unlockers").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	ambWatchdog := service_validity_watchdog.NewWatchTransfersValidity(ambBridge, sideBridge.GetContract())
	sideWatchdog := service_validity_watchdog.NewWatchTransfersValidity(sideBridge, ambBridge.GetContract())
	go service_unlock.NewUnlockTransfers(ambBridge, ambWatchdog).Run()
	go service_unlock.NewUnlockTransfers(sideBridge, sideWatchdog).Run()
}

func runTriggers(cfg *config.Triggers, ambBridge *amb.Bridge, sideBridge networks.Bridge, logger zerolog.Logger) {
	logger.Info().Str("service", "triggers").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	go service_trigger.NewTriggerTransfers(ambBridge).Run()
	go service_trigger.NewTriggerTransfers(sideBridge).Run()
}

func runFeeApi(cfg *config.FeeApi, ambBridge, sideBridge networks.Bridge, logger zerolog.Logger) {
	logger.Info().Str("service", "fee api").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	feeAmb, err := fee_helper.NewFeeHelper(ambBridge, sideBridge, cfg.Amb, cfg.Side)
	if err != nil {
		logger.Fatal().Err(err).Msg("feeAmb not created")
	}
	feeSide, err := fee_helper.NewFeeHelper(sideBridge, ambBridge, cfg.Side, cfg.Amb)
	if err != nil {
		logger.Fatal().Err(err).Msg("feeSide not created")
	}

	feeService := fee.NewFee(feeAmb, feeSide)
	feeApi := &api.FeeAPI{Service: feeService}
	if err = feeApi.Run(cfg.Endpoint, cfg.Ip, cfg.Port, &logger); err != nil {
		logger.Fatal().Err(err).Msg("failed to serve HTTP server (Fee Api endpoint)")
	}
}

func runPrometheus(cfg *config.Prometheus, logger zerolog.Logger) {
	logger.Info().Str("service", "prometheus").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	if err := metric.ServeEndpoint(cfg.Ip, cfg.Port); err != nil {
		logger.Fatal().Err(err).Msg("failed to serve HTTP server (Prometheus endpoint)")
	}
}

func runHealth(addr string, logger zerolog.Logger) {
	// it's actually not /health but * instead :)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(addr, handler); err != nil {
		logger.Fatal().Err(err).Msg("failed to serve HTTP server (Health endpoint)")
	}

}
