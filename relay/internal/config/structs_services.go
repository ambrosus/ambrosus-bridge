package config

type (
	Submitters struct {
		enable    `mapstructure:",squash"`
		AmbToSide SubmitterVariants `mapstructure:"ambToSide"`
		SideToAmb SubmitterVariants `mapstructure:"sideToAmb"`

		AmbFaucet AmbFaucetConfig `mapstructure:"ambFaucet"`
	}

	SubmitterVariants struct {
		Variant string `mapstructure:"variant"`

		Mpc *SubmitterMpc `mapstructure:"mpc"`
	}

	SubmitterMpc struct {
		IsServer    bool     `mapstructure:"isServer"`
		MeID        string   `mapstructure:"meID"`
		PartyIDs    []string `mapstructure:"partyIds"` // for sign operation
		Threshold   int      `mapstructure:"threshold"`
		AccessToken string   `mapstructure:"accessToken"` // client puts this token into request; server checks this token
		ServerURL   string   `mapstructure:"serverURL"`   // client connect to this url; server listen on this url
		SharePath   string   `mapstructure:"sharePath"`
	}

	AmbFaucetConfig struct {
		enable        `mapstructure:",squash"`
		FaucetAddress string `mapstructure:"faucetAddress"`
		MinBalance    int64  `mapstructure:"minBalance"`
		SendAmount    int64  `mapstructure:"sendAmount"`
	}
)

type (
	ValidityWatchdogs struct {
		enable        `mapstructure:",squash"`
		EnableForAmb  bool `mapstructure:"enableForAmb"`
		EnableForSide bool `mapstructure:"enableForSide"`
	}
	PauseUnpauseWatchdogs struct {
		enable `mapstructure:",squash"`
	}
	Triggers struct {
		enable `mapstructure:",squash"`
	}
	Unlockers struct {
		enable `mapstructure:",squash"`
	}
)

type (
	FeeApi struct {
		enable   `mapstructure:",squash"`
		Ip       string `mapstructure:"ip"`
		Port     int    `mapstructure:"port"`
		Endpoint string `mapstructure:"endpoint"`

		Amb  FeeApiNetwork `mapstructure:"amb"`
		Side FeeApiNetwork `mapstructure:"side"`
	}
	FeeApiNetwork struct {
		PrivateKey     string  `mapstructure:"privateKey"`
		MinBridgeFee   float64 `mapstructure:"minBridgeFeeUSD"`
		MinTransferFee float64 `mapstructure:"minTransferFeeUSD"`
		FeeApiUrl      string  `mapstructure:"feeApiUrl"`
	}
)

type (
	enable struct {
		Enable bool `mapstructure:"enable"`
	}
)
