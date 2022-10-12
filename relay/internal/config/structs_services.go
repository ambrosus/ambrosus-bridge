package config

type (
	Submitters struct {
		enable    `mapstructure:",squash"`
		AmbToSide string `mapstructure:"ambToSide"`
		SideToAmb string `mapstructure:"sideToAmb"`

		Aura *SubmitterAura `mapstructure:"aura"`
		Pow  *SubmitterPoW  `mapstructure:"pow"`
		Posa *SubmitterPoSA `mapstructure:"posa"`
	}

	SubmitterAura struct {
		VSContractAddr            string `mapstructure:"vsContractAddr"`
		FinalizeServiceUrl        string `mapstructure:"finalizeServiceUrl"`
		ReceiverBridgeMaxTxSizeKB uint64 `mapstructure:"receiverBridgeMaxTxSizeKB"`
	}
	SubmitterPoW struct {
		EthashDir            string `mapstructure:"ethashDir"`
		EthashKeepPrevEpochs uint64 `mapstructure:"ethashKeepPrevEpochs"`
		EthashGenNextEpochs  uint64 `mapstructure:"ethashGenNextEpochs"`
	}
	SubmitterPoSA struct {
		ReceiverBridgeMaxTxSizeKB uint64 `mapstructure:"receiverBridgeMaxTxSizeKB"`
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

		ExplorerURL          string `mapstructure:"explorerURL"`
		TransferFeeRecipient string `mapstructure:"transferFeeRecipient"`
	}
)

type (
	enable struct {
		Enable bool `mapstructure:"enable"`
	}
)
