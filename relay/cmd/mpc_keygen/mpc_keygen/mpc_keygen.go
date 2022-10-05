package mpc_keygen

import (
	"errors"
	"fmt"
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/rs/zerolog"
)

const ClientSharePath = "share-server.json"
const ServerSharePath = "share-client.json"

func GetAndSaveServerShare(tss *tss_wrap.Mpc, logger *zerolog.Logger) ([]byte, error) {
	return getAndSaveShare(tss, ServerSharePath, logger)
}
func GetAndSaveClientShare(tss *tss_wrap.Mpc, logger *zerolog.Logger) ([]byte, error) {
	return getAndSaveShare(tss, ClientSharePath, logger)
}
func getAndSaveShare(tss *tss_wrap.Mpc, sharePath string, logger *zerolog.Logger) ([]byte, error) {
	logger.Info().Msg("Getting the share...")
	share, err := tss.Share()
	if err != nil {
		return nil, fmt.Errorf("failed to get the share: %w", err)
	}

	logger.Info().Msg("Saving the share...")
	err = os.WriteFile(sharePath, share, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to save the share: %w", err)
	}
	logger.Info().Msg("The share has been saved!")

	return share, nil
}

func IsServerShareExist() bool {
	return isFileExist(ServerSharePath)
}
func IsClientShareExist() bool {
	return isFileExist(ClientSharePath)
}

func ReadServerShare() ([]byte, error) {
	return os.ReadFile(ServerSharePath)
}
func ReadClientShare() ([]byte, error) {
	return os.ReadFile(ClientSharePath)
}

func isFileExist(filepath string) bool {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
