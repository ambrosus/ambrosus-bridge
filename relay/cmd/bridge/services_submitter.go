package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/untrustless2"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/untrustless_mpc"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/client"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/server"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
)

func createSubmitter(cfg *config.SubmitterVariants, submitterBridge networks.Bridge, receiverBridge service_submit.Receiver) (service_submit.Submitter, error) {
	switch cfg.Variant {
	case "":
		return nil, nil
	case "untrustless2":
		return untrustless2.NewSubmitterUntrustless(submitterBridge, &untrustless2.ReceiverUntrustless2{Receiver: receiverBridge})
	case "untrustless-mpc":
		return submitterMpc(cfg.Mpc, submitterBridge, receiverBridge)
	default:
		return nil, fmt.Errorf("unknown submitter variant: %s", cfg.Variant)
	}
}

func submitterMpc(cfg *config.SubmitterMpc, submitterBridge networks.Bridge, receiverBridge service_submit.Receiver) (service_submit.Submitter, error) {
	logger := submitterBridge.GetLogger()

	share, err := os.ReadFile(cfg.SharePath)
	if err != nil {
		return nil, fmt.Errorf("can't read share: %w", err)
	}
	mpcc, err := tss_wrap.NewMpcWithShare(cfg.MeID, cfg.Threshold, share, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create MPC client: %w", err)
	}

	var mpcSigner untrustless_mpc.MpcSigner
	if cfg.IsServer {
		server_ := server.NewServer(mpcc, cfg.AccessToken, logger)
		go http.ListenAndServe(cfg.ServerURL, server_)
		mpcSigner = server_
	} else {
		client_ := client.NewClient(mpcc, cfg.ServerURL, cfg.AccessToken, logger)
		mpcSigner = client_
	}

	receiver, err := untrustless_mpc.NewReceiverUntrustlessMpc(receiverBridge, mpcSigner, cfg.PartyIDs)
	if err != nil {
		return nil, fmt.Errorf("NewSubmitterUntrustlessMpc: %w", err)
	}
	submitter, err := untrustless_mpc.NewSubmitterUntrustlessMpc(submitterBridge, receiver, cfg.IsServer)
	if err != nil {
		return nil, fmt.Errorf("NewSubmitterUntrustlessMpc: %w", err)
	}

	return submitter, nil
}
