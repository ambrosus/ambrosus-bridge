package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/ambrosus/ambrosus-bridge/relay/cmd"
	"github.com/ambrosus/ambrosus-bridge/relay/cmd/mpc_keygen/mpc_keygen"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/client"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/server"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
)

func main() {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error initialize config")
	}

	logger := cmd.CreateLogger(cfg.ExtLoggers)
	cfgSubmitters := cfg.Submitters

	switch cfg.Submitters.AmbToSide {
	case "untrustless-mpc-server":
		mpcc := tss_wrap.NewMpc(&tss_wrap.MpcConfig{MeID: int(cfgSubmitters.MpcServer.MeID), PartyLen: int(cfgSubmitters.MpcServer.PartyLen), Threshold: int(cfgSubmitters.MpcServer.Threshold)}, &logger)
		server_ := server.NewServer(mpcc, &logger)
		go http.ListenAndServe(fmt.Sprintf(":%d", cfg.Submitters.MpcServer.Port), server_)
		go server_.Run()

		err := server_.Keygen()
		if err != nil {
			logger.Fatal().Err(err).Msg("error on untrustless mpc server keygen")
		}

		_, err = mpc_keygen.GetAndSaveServerShare(server_.Tss, &logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("error on getting and saving share")
		}

	case "untrustless-mpc-client":
		mpcc := tss_wrap.NewMpc(&tss_wrap.MpcConfig{MeID: int(cfgSubmitters.MpcClient.MeID), PartyLen: int(cfgSubmitters.MpcClient.PartyLen), Threshold: int(cfgSubmitters.MpcClient.Threshold)}, &logger)
		client := client.NewClient(mpcc, cfgSubmitters.MpcClient.ServerURL, &logger)

		err := client.Keygen()
		if err != nil {
			logger.Err(err).Msg("error on untrustless mpc client keygen")
		}

		_, err = mpc_keygen.GetAndSaveClientShare(client.Tss, &logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("error on getting and saving share")
		}
	}
}
