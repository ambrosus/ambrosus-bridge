package config

type (
	Submitters struct {
		enable    `mapstructure:",squash"`
		AmbToSide string `mapstructure:"ambToSide"`
		SideToAmb string `mapstructure:"sideToAmb"`

		Aura      *SubmitterAura      `mapstructure:"aura"`
		Pow       *SubmitterPoW       `mapstructure:"pow"`
		Posa      *SubmitterPoSA      `mapstructure:"posa"`
		MpcClient *SubmitterMpcClient `mapstructure:"mpcClient"`
		MpcServer *SubmitterMpcServer `mapstructure:"mpcServer"`
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
	SubmitterMpcClient struct {
		MeID      uint64 `mapstructure:"meID"`
		PartyLen  uint64 `mapstructure:"partyLen"`
		Threshold uint64 `mapstructure:"threshold"`
		ServerURL string `mapstructure:"serverURL"`
	}
	SubmitterMpcServer struct {
		MeID      uint64 `mapstructure:"meID"`
		PartyLen  uint64 `mapstructure:"partyLen"`
		Threshold uint64 `mapstructure:"threshold"`
		Port      uint64 `mapstructure:"port"`
	}
)

type (
	Watchdogs struct {
		enable        `mapstructure:",squash"`
		EnableForAmb  bool `mapstructure:"enableForAmb"`
		EnableForSide bool `mapstructure:"enableForSide"`
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
