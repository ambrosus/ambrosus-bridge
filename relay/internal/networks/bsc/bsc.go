package bsc

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/limit_filtering"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
)

const BridgeName = "binance"

type Bridge struct {
	nc.CommonBridge
}

// New creates a new ethereum bridge.
func New(cfg *config.Network, sideBridgeName string, baseLogger zerolog.Logger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}

	commonBridge.Logger = baseLogger.With().Str("bridge", BridgeName).Str("sideBridge", sideBridgeName).Logger()

	// ///////////////////
	origin := nc.GetAmbrosusOrigin()
	filterLogsFromBlock, err := strconv.Atoi(os.Getenv("BSC_FILTER_LOGS_FROM_BLOCK"))
	if err != nil {
		return nil, fmt.Errorf("invalid BSC_FILTER_LOGS_FROM_BLOCK: %w", err)
	}
	commonBridge.Logger.Info().Msgf("Set BSC_FILTER_LOGS_FROM_BLOCK to %d", filterLogsFromBlock)

	filterLogsLimitBlocks, err := strconv.Atoi(os.Getenv("BSC_FILTER_LOGS_LIMIT_BLOCKS"))
	if err != nil {
		return nil, fmt.Errorf("invalid BSC_FILTER_LOGS_LIMIT_BLOCKS: %w", err)
	}
	commonBridge.Logger.Info().Msgf("Set BSC_FILTER_LOGS_LIMIT_BLOCKS to %d", filterLogsLimitBlocks)

	rpcHTTPClient, err := rpc.DialHTTP(cfg.HttpURL)
	if err != nil {
		return nil, fmt.Errorf("dial http: %w", err)
	}
	rpcHTTPClient.SetHeader("Origin", origin)

	limitFilteringClient := limit_filtering.NewClient(rpcHTTPClient, int64(filterLogsFromBlock), int64(filterLogsLimitBlocks))
	commonBridge.Client = limitFilteringClient

	// Creating a new bridge contract instance.
	commonBridge.Contract, err = bindings.NewBridge(commonBridge.ContractAddress, commonBridge.Client)
	if err != nil {
		return nil, fmt.Errorf("create contract http: %w", err)
	}

	// Create websocket instances if wsUrl provided
	if cfg.WsURL != "" {
		rpcWSClient, err := rpc.DialWebsocket(context.Background(), cfg.WsURL, origin)
		if err != nil {
			return nil, fmt.Errorf("dial ws: %w", err)
		}
		commonBridge.WsClient = limit_filtering.NewClient(rpcWSClient, int64(filterLogsFromBlock), int64(filterLogsLimitBlocks))

		commonBridge.WsContract, err = bindings.NewBridge(commonBridge.ContractAddress, commonBridge.WsClient)
		if err != nil {
			return nil, fmt.Errorf("create contract ws: %w", err)
		}
	}

	return &Bridge{
		CommonBridge: commonBridge,
	}, nil
}
