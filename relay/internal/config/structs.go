package config

type (
	Network struct {
		HttpURL      string `mapstructure:"httpUrl"`
		WsURL        string `mapstructure:"wsUrl"`
		ContractAddr string `mapstructure:"contractAddr"`
		PrivateKey   string `mapstructure:"privateKey"`
	}

	AMBConfig struct {
		Network        `mapstructure:",squash"`
		VSContractAddr string `mapstructure:"vsContractAddr"`
	}

	BSCConfig struct {
		Network `mapstructure:",squash"`
	}

	ETHConfig struct {
		Network              `mapstructure:",squash"`
		EthashDir            string `mapstructure:"ethashDir"`
		EthashKeepPrevEpochs uint64 `mapstructure:"ethashKeepPrevEpochs"`
		EthashGenNextEpochs  uint64 `mapstructure:"ethashGenNextEpochs"`
	}
)

type (
	Config struct {
		Networks   Networks        `mapstructure:"network"`
		ExtLoggers ExternalLoggers `mapstructure:"externalLogger"`
		Prometheus Prometheus      `mapstructure:"prometheus"`
		FeeApi     FeeApi          `mapstructure:"feeApi"`
		IsRelay    bool            `mapstructure:"isRelay"`
		IsWatchdog bool            `mapstructure:"isWatchdog"`
	}

	Networks struct {
		AMB *AMBConfig `mapstructure:"amb"`
		ETH *ETHConfig `mapstructure:"eth"`
		BSC *BSCConfig `mapstructure:"bsc"`
	}

	ExternalLoggers struct {
		Telegram TelegramLogger `mapstructure:"telegram"`
	}
	TelegramLogger struct {
		Enable bool   `mapstructure:"enable"`
		Token  string `mapstructure:"token"`
		ChatId int    `mapstructure:"chatId"`
	}

	FeeApi struct {
		Enable   bool   `mapstructure:"enable"`
		Ip       string `mapstructure:"ip"`
		Port     int    `mapstructure:"port"`
		Endpoint string `mapstructure:"endpoint"`

		Amb  FeeApiNetwork `mapstructure:"amb"`
		Side FeeApiNetwork `mapstructure:"side"`
	}
	FeeApiNetwork struct {
		PrivateKey         string  `mapstructure:"privateKey"`
		MinBridgeFee       float64 `mapstructure:"minBridgeFeeUSD"`
		DefaultTransferFee int64   `mapstructure:"defaultTransferFee"`
	}

	Prometheus struct {
		Enable bool   `mapstructure:"enable"`
		Ip     string `mapstructure:"ip"`
		Port   int    `mapstructure:"port"`
	}
)
