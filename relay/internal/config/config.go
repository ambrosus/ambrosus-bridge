package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const defaultConfigPath string = "configs/main"

func LoadConfig() (*Config, error) {
	log.Debug().Msg("Loading config...")

	dir, file := filepath.Split(getConfigPath())
	viper.AddConfigPath(dir)
	viper.SetConfigName(file)
	viper.SetConfigType("json")

	// ex: use `NETWORK_AMB_PRIVATEKEY` env var to override `network.amb.privateKey` key
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read cfg: %w", err)
	}

	cfg := new(Config)
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal cfg: %w", err)
	}

	return cfg, nil
}

func getConfigPath() string {
	stage := os.Getenv("STAGE")     // dev / test / main / ...
	network := os.Getenv("NETWORK") // eth / bsc / ...

	if stage != "" && network != "" {
		configPath := fmt.Sprintf("configs/%s-%s", stage, network)
		log.Info().Msgf("`STAGE` and `NETWORK` env var set, using config path: %v", configPath)
		return configPath
	} else if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		log.Info().Msgf("`STAGE` or `NETWORK` env var not set, using `CONFIG_PATH` env var: %v", configPath)
		return configPath
	} else {
		log.Info().Msgf("(`STAGE` or `NETWORK`) and `CONFIG_PATH` env var not set, using default config path: %v", defaultConfigPath)
		return defaultConfigPath
	}
}
