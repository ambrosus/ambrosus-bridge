package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const defaultConfigPath string = "configs/main"

var ErrPrivateKeyNotFound = errors.New("private key not found in environment")

type (
	Config struct {
		AMB        AMBConfig
		ETH        ETHConfig
		Telegram   TelegramLogger
		Prometheus Prometheus
	}

	Network struct {
		HttpURL      string `mapstructure:"httpUrl"`
		WsURL        string `mapstructure:"wsUrl"`
		ContractAddr string `mapstructure:"contractAddr"`
		PrivateKey   string
	}

	AMBConfig struct {
		Network
		VSContractAddr string `mapstructure:"vsContractAddr"`
	}

	ETHConfig struct {
		Network
		EpochLength uint64 `mapstructure:"epochLength"`
		EthashDir   string `mapstructure:"ethashDir"`
	}

	TelegramLogger struct {
		Enable bool   `mapstructure:"enable"`
		Token  string `mapstructure:"token"`
		ChatId int    `mapstructure:"chat-id"`
	}

	Prometheus struct {
		Enable bool   `mapstructure:"enable"`
		Ip     string `mapstructure:"ip"`
		Port   int    `mapstructure:"port"`
	}
)

func Init() (*Config, error) {
	log.Debug().Msg("Initialize config...")

	if err := parseConfigFile(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := setFromEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseConfigFile() error {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = defaultConfigPath
	}

	log.Debug().Msgf("Parsing config file: %s", configPath)

	dir, file := filepath.Split(configPath)

	viper.SetConfigType("json")
	viper.AddConfigPath(dir)
	viper.SetConfigName(file)

	return viper.ReadInConfig()

}

func unmarshal(cfg *Config) error {
	log.Debug().Msg("Unmarshal config keys...")

	if err := viper.UnmarshalKey("network.amb", &cfg.AMB); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("network.amb", &cfg.AMB.Network); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("network.eth", &cfg.ETH); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("externalLogger.telegram", &cfg.Telegram); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("prometheus", &cfg.Prometheus); err != nil {
		return err
	}

	return viper.UnmarshalKey("network.eth", &cfg.ETH.Network)
}

func setFromEnv(cfg *Config) error {
	log.Debug().Msg("Set from environment configurations...")

	cfg.AMB.PrivateKey = os.Getenv("AMB_PRIVATE_KEY")
	if cfg.AMB.PrivateKey == "" {
		return ErrPrivateKeyNotFound
	}
	cfg.ETH.PrivateKey = os.Getenv("ETH_PRIVATE_KEY")
	if cfg.ETH.PrivateKey == "" {
		return ErrPrivateKeyNotFound
	}
	return nil
}
