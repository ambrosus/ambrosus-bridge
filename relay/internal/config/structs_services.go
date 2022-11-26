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

		Aura *SubmitterAura `mapstructure:"aura"`
		Pow  *SubmitterPoW  `mapstructure:"pow"`
		Posa *SubmitterPoSA `mapstructure:"posa"`
		Mpc  *SubmitterMpc  `mapstructure:"mpc"`
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
		enable     `mapstructure:",squash"`
		PrivateKey string `mapstructure:"privateKey"`
		MinBalance int64  `mapstructure:"minBalance"`
		SendAmount int64  `mapstructure:"sendAmount"`
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

		ExplorerURL                         string   `mapstructure:"explorerURL"`
		TransferFeeRecipient                string   `mapstructure:"transferFeeRecipient"`
		TransferFeeIncludedTxsFromAddresses []string `mapstructure:"transferFeeIncludedTxsFromAddresses"`
		TransferFeeTxsFromBlock             uint64   `mapstructure:"transferFeeTxsFromBlock"`
	}
)

type (
	enable struct {
		Enable bool `mapstructure:"enable"`
	}
)
