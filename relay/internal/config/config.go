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

func LoadDefaultConfig() (*Config, error) {
	log.Debug().Msg("Loading config...")

	v, err := LoadConfig(getConfigPath())
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	if err = v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal cfg: %w", err)
	}

	return cfg, nil
}

func LoadConfig(path string) (*viper.Viper, error) {
	dir, file := filepath.Split(path)
	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigName(file)
	v.SetConfigType("json")

	// ex: use `NETWORK_AMB_PRIVATEKEY` env var to override `network.amb.privateKey` key
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read cfg: %w", err)
	}
	return v, nil
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
