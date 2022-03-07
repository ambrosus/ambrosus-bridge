package config

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const defaultConfigPath string = "configs/main"

var ErrPrivateKeyNotFound = errors.New("private key not found in environment")

type (
	Config struct {
		AMB      AMBConfig
		ETH      ETHConfig
		Telegram TelegramLogger
	}

	Network struct {
		URL          string `mapstructure:"url"`
		ContractAddr string `mapstructure:"contract-addr"`
		PrivateKey   *ecdsa.PrivateKey
	}

	AMBConfig struct {
		Network
		VSContractAddr string `mapstructure:"vs-contract-addr"`
	}

	ETHConfig struct {
		Network
		EthashPath  string `mapstructure:"ethash-path"`
		EpochLenght uint64 `mapstructure:"epoch-lenght"`
	}

	TelegramLogger struct {
		Token  string `mapstructure:"token"`
		ChatId int    `mapstructure:"chat-id"`
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

	if err := viper.UnmarshalKey("external-logger.telegram", &cfg.Telegram); err != nil {
		return err
	}

	return viper.UnmarshalKey("network.eth", &cfg.ETH.Network)
}

func setFromEnv(cfg *Config) error {
	log.Debug().Msg("Set from environment configurations...")

	ambPrivateKey, err := parsePK(os.Getenv("AMB_PRIVATE_KEY"))
	if err != nil {
		return err
	}
	ethPrivateKey, err := parsePK(os.Getenv("ETH_PRIVATE_KEY"))
	if err != nil {
		return err
	}

	cfg.AMB.PrivateKey = ambPrivateKey
	cfg.ETH.PrivateKey = ethPrivateKey

	return nil
}

func parsePK(pk string) (*ecdsa.PrivateKey, error) {
	if pk == "" {
		return nil, ErrPrivateKeyNotFound
	}

	b, err := hex.DecodeString(pk)
	if err != nil {
		return nil, err
	}
	p, err := crypto.ToECDSA(b)
	if err != nil {
		return nil, err
	}

	return p, nil
}
