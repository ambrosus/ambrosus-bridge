package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"

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

	logger := log.Logger
	mpcc := tss_wrap.NewMpc(cfg.Submitters.Mpc.MeID, cfg.Submitters.Mpc.PartyLen, &logger)

	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	if cfg.Submitters.Mpc.IsServer {
		server_ := server.NewServer(mpcc, &logger)
		go http.ListenAndServe(cfg.Submitters.Mpc.ServerURL, server_)

		err := server_.Keygen(ctx)
		if err != nil {
			logger.Fatal().Err(err).Msg("error on untrustless mpc server keygen")
		}
	} else {
		client_ := client.NewClient(mpcc, cfg.Submitters.Mpc.ServerURL, nil, &logger)

		err := client_.Keygen(ctx)
		if err != nil {
			logger.Err(err).Msg("error on untrustless mpc client keygen")
		}
	}
	saveShare(mpcc)

	logger.Info().Msg("mpc keygen finished")
}

func saveShare(tss *tss_wrap.Mpc) {
	share, err := tss.Share()
	if err != nil {
		panic(err)
	}

	sharePath := os.Getenv("SHARE_PATH")
	err = os.WriteFile(sharePath, share, 0644)
	if err != nil {
		panic(err)
	}
}
