package main

import (
	"context"
	"fmt"
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

	if !isShareExist(cfg.Submitters.AmbToSide.Mpc.SharePath) {
		fmt.Println("AmbToSide keygen begin")
		keygen(cfg.Submitters.AmbToSide.Mpc)
		fmt.Println("AmbToSide keygen done")
	} else {
		fmt.Println("AmbToSide share already exist")
	}

	if !isShareExist(cfg.Submitters.SideToAmb.Mpc.SharePath) {
		fmt.Println("SideToAmb keygen begin")
		keygen(cfg.Submitters.SideToAmb.Mpc)
		fmt.Println("SideToAmb keygen done")
	} else {
		fmt.Println("SideToAmb share already exist")
	}

}

func keygen(cfg *config.SubmitterMpc) {
	logger := log.Logger
	mpcc := tss_wrap.NewMpc(cfg.MeID, cfg.PartyLen, &logger)

	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	if cfg.IsServer {
		server_ := server.NewServer(mpcc, &logger)
		go http.ListenAndServe(cfg.ServerURL, server_)

		err := server_.Keygen(ctx)
		if err != nil {
			panic(err)
		}
	} else {
		client_ := client.NewClient(mpcc, cfg.ServerURL, nil, &logger)

		err := client_.Keygen(ctx)
		if err != nil {
			panic(err)
		}
	}
	saveShare(mpcc, cfg.SharePath)
}

func saveShare(tss *tss_wrap.Mpc, sharePath string) {
	share, err := tss.Share()
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(sharePath, share, 0644)
	if err != nil {
		panic(err)
	}
}

func isShareExist(sharePath string) bool {
	_, err := os.Stat(sharePath)
	return err == nil
}
