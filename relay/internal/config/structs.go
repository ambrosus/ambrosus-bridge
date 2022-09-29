package config

type (
	Config struct {
		Networks *Networks `mapstructure:"networks"`

		Submitters        *Submitters        `mapstructure:"submitters"`
		ValidityWatchdogs *ValidityWatchdogs `mapstructure:"validityWatchdogs"`
		Triggers          *Triggers          `mapstructure:"triggers"`
		Unlockers         *Unlockers         `mapstructure:"unlockers"`
		FeeApi            *FeeApi            `mapstructure:"feeApi"`

		ExtLoggers *ExternalLoggers `mapstructure:"externalLogger"`
		Prometheus *Prometheus      `mapstructure:"prometheus"`
	}
)

type (
	Networks struct {
		SideBridgeNetwork string              `mapstructure:"sideBridgeNetwork"`
		Networks          map[string]*Network `mapstructure:",remain"`
	}
	Network struct {
		HttpURL      string `mapstructure:"httpUrl"`
		WsURL        string `mapstructure:"wsUrl"`
		ContractAddr string `mapstructure:"contractAddr"`
		PrivateKey   string `mapstructure:"privateKey"`
	}
)

type (
	ExternalLoggers struct {
		Telegram TelegramLogger `mapstructure:"telegram"`
	}
	TelegramLogger struct {
		enable `mapstructure:",squash"`
		Token  string `mapstructure:"token"`
		ChatId int    `mapstructure:"chatId"`
	}
)

type (
	Prometheus struct {
		enable `mapstructure:",squash"`
		Ip     string `mapstructure:"ip"`
		Port   int    `mapstructure:"port"`
	}
)
