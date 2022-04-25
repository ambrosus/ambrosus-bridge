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

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Info().Msgf("`CONFIG_PATH` env var not set, using default config path: %v", defaultConfigPath)
		configPath = defaultConfigPath
	}
	dir, file := filepath.Split(configPath)

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
