const Rory = "0x40B7d71E70fA6311cB0b300c1Ba6926A2A9000b8";
const Lang = "0xEB1c6a8a84063B1cef8B9a23AB87Bf926035A21a";
const Kevin = "0x260cfE305cA40CaE1a32Ba7611137eF4d7146233";
const Andrey = "0xb017DcCC473499C83f1b553bE564f3CeAf002254";

const Master_BSC_Amb = "0xDE99F453d76D5f485d90244319593b12D5eeDe8B";
const Master_BSC_Bsc = "0x64C17b651d0c7d2dF92A7f6E7707C27C6e0535D8";
const Master_ETH_Amb = "0x0f071e1785e3E115360E04c9C8D53e958E6f85FE";
const Master_ETH_Eth = "0xFAE075e12116FBfE65c58e1Ef0E6CA959cA37ded";

const AdminAmb = "0x203F8dCBce61d9BFa48A95F588152Da073d9c693";
const AdminBsc = "0x3602541a2015A7867532EEF7e3A4eB2ED436a715";
const AdminEth = "0xb6db3c082d25C0dC26602334B5F3D270091DC422";


const multisig = {
  "admins": [
    Andrey,
    "0xf3cb1e06e72233e0350F7588D11a2585B268C7AD", // Andrey R
    "0xe620e1F969Bc3a24Ac96D527220AD6B6e2d12843", // Olena
    "0xe8d204d3b12e643888debd525e2f8034fc7eb855", // Alina
    // "0xFe231E979dAe1D87Cd961Aa353221C0c754941d8", // Oggy OLD
    "0x5700F8e0ae3d80964f7718EA625E3a2CB4D2096d", // Oggy
    "0x59C66c7E3D3F239961c87A4627F86a5a9049407C", // Svin
    "0xa5E32D3fB342D9Ed3135fD5cb59a102AC8ED7B85" // Valar
  ],
  "threshold": 4
}


const prod_addresses = {

  // amb part of amb-bsc bridge
  "BSC_AmbBridge": {
    "adminAddress": AdminAmb,
    "relayAddress": Master_BSC_Amb,
    "masterRelayAddress": Master_BSC_Amb,
    "feeProviderAddress": Master_BSC_Amb,
    "watchdogsAddresses": [
      Kevin, Lang, Rory, Andrey,
      Master_BSC_Amb,
      AdminAmb
    ],
    "transferFeeRecipient": "0xba971570E4352a700de3ca57Fe882E2d4C70F42F",
    "bridgeFeeRecipient": "0x893dF80919D7e89A56ED2668466ad4EB84C367a0",
    multisig,
  },

  // bsc part of amb-bsc bridge
  "BSC_BscBridge": {
    "adminAddress": AdminBsc,
    "relayAddress": Master_BSC_Bsc,
    "masterRelayAddress": Master_BSC_Bsc,
    "feeProviderAddress": Master_BSC_Bsc,
    "watchdogsAddresses": [
      Kevin, Lang, Rory, Andrey,
      Master_BSC_Bsc,
      AdminBsc
    ],
    "transferFeeRecipient": "0x90268508333A15BC932f63faD95E40B7a251eFB1",
    "bridgeFeeRecipient": "0x4CD40e00eBD87acbA000071439B647Aaa2810683",
    multisig,
  },

  // amb part of amb-eth bridge
  "ETH_AmbBridge": {
    "adminAddress": AdminAmb,
    "relayAddress": "0x3BB6161dc6E5380831833b223D467b715e05F16E",
    "masterRelayAddress": Master_ETH_Amb,
    "feeProviderAddress": Master_ETH_Amb,
    "watchdogsAddresses": [
      Kevin, Lang, Rory, Andrey,
      Master_ETH_Amb,
      AdminAmb
    ],
    "transferFeeRecipient": "0xba971570E4352a700de3ca57Fe882E2d4C70F42F",
    "bridgeFeeRecipient": "0xD3bd2Ac57e3FE109BF29dF67b251E1BfC96DA7b8",
    multisig,
  },

  // eth part of amb-eth bridge
  "ETH_EthBridge": {
    "adminAddress": AdminEth,
    "relayAddress": "0x3BB6161dc6E5380831833b223D467b715e05F16E",
    "masterRelayAddress": Master_ETH_Eth,
    "feeProviderAddress": Master_ETH_Eth,
    "watchdogsAddresses": [
      Kevin, Lang, Rory, Andrey,
      Master_ETH_Eth,
      AdminEth
    ],
    "transferFeeRecipient": "0x0BD1dD0F45f0F180001e6BD84477B985e8Ce7295",
    "bridgeFeeRecipient": "0x3f3DA18920Ca95bC55353E8e4Fd3C1123D3A07B4",
    multisig,
  }
}

export default prod_addresses;
