package config

type (
	Submitters struct {
		enable    `mapstructure:",squash"`
		AmbToSide bool `mapstructure:"ambToSide"`
		SideToAmb bool `mapstructure:"sideToAmb"`

		Aura *SubmitterAura `mapstructure:"aura"`
		Pow  *SubmitterPoW  `mapstructure:"pow"`
		//Posa    *SubmitterPoSA `mapstructure:"posa"`  PoSA doesn't need any cfg
	}

	SubmitterAura struct {
		VSContractAddr     string `mapstructure:"vsContractAddr"`
		FinalizeServiceUrl string `mapstructure:"finalizeServiceUrl"`
	}
	SubmitterPoW struct {
		EthashDir            string `mapstructure:"ethashDir"`
		EthashKeepPrevEpochs uint64 `mapstructure:"ethashKeepPrevEpochs"`
		EthashGenNextEpochs  uint64 `mapstructure:"ethashGenNextEpochs"`
	}
)

type (
	Watchdogs struct {
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
		PrivateKey         string  `mapstructure:"privateKey"`
		MinBridgeFee       float64 `mapstructure:"minBridgeFeeUSD"`
		DefaultTransferFee float64 `mapstructure:"defaultTransferFeeWei"`
	}
)

type (
	enable struct {
		Enable bool `mapstructure:"enable"`
	}
)
