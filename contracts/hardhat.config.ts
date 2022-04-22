import "@nomiclabs/hardhat-waffle";
import "hardhat-deploy";
import "hardhat-gas-reporter";
import "solidity-coverage";
import "hardhat-abi-exporter";


import * as dotenv from "dotenv";

dotenv.config();
// todo add other roles
const PK = [process.env.PRIVATEKEY || "00000000"];

module.exports = {
  networks: {
    hardhat: {
      blockGasLimit: 40000000, // amb value
      hardfork: "byzantium",
      companionNetworks: {amb: 'hardhat'},
    },

    "dev/eth": {
      url: "https://ropsten.infura.io/v3/" + process.env.INFURA_KEY,
      accounts: PK,
      tags: ["eth", "devnet"],
      companionNetworks: {amb: 'amb_dev'},
    },
    "test/eth": {
      url: "https://ropsten.infura.io/v3/" + process.env.INFURA_KEY,
      accounts: PK,
      tags: ["eth", "testnet"],
      companionNetworks: {amb: 'amb_test'},
    },
    "main/eth": {
      url: "https://mainnet.infura.io/v3/" + process.env.INFURA_KEY,
      accounts: PK,
      tags: ["eth", "mainnet"],
      companionNetworks: {amb: 'amb_main'},
    },
    "integr/eth": {
      url: "http://127.0.0.1:8502",
      accounts: ["0x51d098d8aee092622149d8f3a79cc7b1ce36ff97fadaa2fbd623c65badeefadc"],
      tags: ["eth", "integr"],
      companionNetworks: {amb: 'amb_integr'},
    },

    "dev/amb": {
      url: "https://network.ambrosus-dev.io",
      accounts: PK,
      tags: ["amb", "devnet"],
      hardfork: "byzantium",
    },
    "test/amb": {
      url: "https://network.ambrosus-test.io",
      accounts: PK,
      tags: ["amb", "testnet"],
      hardfork: "byzantium",
    },
    "main/amb": {
      url: "https://network.ambrosus.io",
      accounts: PK,
      tags: ["amb", "mainnet"],
      hardfork: "byzantium",
    },
    "integr/amb": {
      url: "http://127.0.0.1:8545",
      accounts: ["0x80f702eb861f36fe8fbbe1a7ccceb04ef7ddef714604010501a5f67c8065d446"],
      tags: ["amb", "integr"],
      hardfork: "byzantium",
    },
  },
  namedAccounts: {
    owner: 0,
    admin: 1,
    relay: 2,
    bridge: 3,
    user: 4,
  },
  etherscan: {
    apiKey: "DY4Z86MQ2D9E24C6HB98PTA79EKJ5TQIFX",
  },
  solidity: {
    compilers: [
      {
        version: "0.8.6",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200,  // todo bigger
          },
          // Note: for amb deploy
          evmVersion: "byzantium"
        },
      }, {
        version: "0.4.22",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200,
          },
          // Note: for amb deploy
          evmVersion: "byzantium"
        },
      },
    ],
  },
  abiExporter: {
    runOnCompile: true,
    clear: true,
    flat: true,
    only: [
      "AmbBridge",
      "EthBridge",
      "ValidatorSet$",
      "sAMB",
      "IWrapper",
    ]
  }
};
