package main

import (
	"fmt"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/cmd/mpc_keygen/mpc_keygen"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/api"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_helper"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/aura"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/mpc"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/posa"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/untrustless"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/untrustless2"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_trigger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_unlock"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_watchdog"
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
	case "untrustless-mpc-server":
		mpcc := tss_wrap.NewMpc(&tss_wrap.MpcConfig{MeID: int(cfg.MpcServer.MeID), PartyLen: int(cfg.MpcServer.PartyLen), Threshold: int(cfg.MpcServer.Threshold)}, &logger)
		server_ := server.NewServer(mpcc, &logger)
		go http.ListenAndServe(fmt.Sprintf(":%d", cfg.MpcServer.Port), server_)

		var share []byte
		if mpc_keygen.IsServerShareExist() {
			share, err = mpc_keygen.ReadServerShare()
			if err != nil {
				logger.Fatal().Err(err).Msg("ambBridgeSubmitter don't created: can't read server's share")
			}
		} else {
			err = server_.Keygen()
			if err != nil {
				logger.Fatal().Err(err).Msg("ambBridgeSubmitter don't created: can't do keygen")
			}

			share, err = mpc_keygen.GetAndSaveServerShare(server_.Tss, &logger)
			if err != nil {
				logger.Fatal().Err(err).Msg("ambBridgeSubmitter don't created: can't get and save server's share")
			}
		}

		err = server_.Tss.SetShare(share)
		if err != nil {
			logger.Fatal().Err(err).Msg("ambBridgeSubmitter don't created: can't set server's share")
		}

		receiver, err := mpc.NewMpcReceiver(sideBridge, server_)
		if err == nil {
			ambSubmitter, err = untrustless2.NewSubmitterUntrustless(ambBridge, &untrustless2.ReceiverUntrustless2{Receiver: receiver})
		}
	case "untrustless-mpc-client":
		mpcc := tss_wrap.NewMpc(&tss_wrap.MpcConfig{MeID: int(cfg.MpcClient.MeID), PartyLen: int(cfg.MpcClient.PartyLen), Threshold: int(cfg.MpcClient.Threshold)}, &logger)
		client := client.NewClient(mpcc, cfg.MpcClient.ServerURL, &logger)

		var share []byte
		if mpc_keygen.IsClientShareExist() {
			share, err = mpc_keygen.ReadClientShare()
			if err != nil {
				logger.Fatal().Err(err).Msg("ambBridgeSubmitter don't created: can't read client's share")
			}
		} else {
			err = client.Keygen()
			if err != nil {
				logger.Fatal().Err(err).Msg("ambBridgeSubmitter don't created: can't do keygen")
			}

			share, err = mpc_keygen.GetAndSaveClientShare(client.Tss, &logger)
			if err != nil {
				logger.Fatal().Err(err).Msg("ambBridgeSubmitter don't created: can't get and save client's share")
			}
		}

		err = client.Tss.SetShare(share)
		if err != nil {
			logger.Fatal().Err(err).Msg("ambBridgeSubmitter don't created: can't set client's share")
		}

		receiver, err := mpc.NewMpcReceiver(sideBridge, client)
		if err == nil {
			ambSubmitter, err = untrustless2.NewSubmitterUntrustless(ambBridge, &untrustless2.ReceiverUntrustless2{Receiver: receiver})
		}
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

func runWatchdogs(cfg *config.Watchdogs, ambBridge *amb.Bridge, sideBridge networks.Bridge, logger zerolog.Logger) {
	logger.Info().Str("service", "watchdog").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	if cfg.EnableForAmb {
		go service_watchdog.NewWatchTransfersValidity(ambBridge, sideBridge.GetContract()).Run()
	}

	if cfg.EnableForSide {
		go service_watchdog.NewWatchTransfersValidity(sideBridge, ambBridge.GetContract()).Run()
	}

}

func runUnlockers(cfg *config.Unlockers, ambBridge *amb.Bridge, sideBridge networks.Bridge, logger zerolog.Logger) {
	logger.Info().Str("service", "watchdog").Bool("enabled", cfg.Enable).Send()
	if !cfg.Enable {
		return
	}

	ambWatchdog := service_watchdog.NewWatchTransfersValidity(ambBridge, sideBridge.GetContract())
	sideWatchdog := service_watchdog.NewWatchTransfersValidity(sideBridge, ambBridge.GetContract())
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
