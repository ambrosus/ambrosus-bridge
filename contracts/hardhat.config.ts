import "@nomiclabs/hardhat-waffle";
import "hardhat-deploy";
import "hardhat-gas-reporter";
import "solidity-coverage";
import "hardhat-abi-exporter";
import * as dotenv from "dotenv";
import {HardhatUserConfig} from "hardhat/types";

dotenv.config();
// todo add other roles
const PK = [
  process.env.PRIVATEKEY || "00000000",
  process.env.SECONDPRIVATEKEY || "00000000"
];

const config: HardhatUserConfig = {

  networks: {
    hardhat: {
      blockGasLimit: 40000000, // amb value
      hardfork: "byzantium",
      companionNetworks: {amb: 'hardhat'},
      initialDate: "2022-05-13T19:03:35"
    },

    "dev/eth": {
      url: "https://ropsten.infura.io/v3/" + process.env.INFURA_KEY,
      accounts: PK,
      tags: ["eth", "devnet"],
      companionNetworks: {amb: 'dev/amb'},
    },
    "test/eth": {
      url: "https://ropsten.infura.io/v3/" + process.env.INFURA_KEY,
      accounts: PK,
      tags: ["eth", "testnet"],
      companionNetworks: {amb: 'test/amb'},
    },
    "main/eth": {
      url: "https://mainnet.infura.io/v3/" + process.env.INFURA_KEY,
      accounts: PK,
      tags: ["eth", "mainnet"],
      companionNetworks: {amb: 'main/amb'},
    },
    "integr/eth": {
      url: "http://127.0.0.1:8502",
      accounts: ["0x51d098d8aee092622149d8f3a79cc7b1ce36ff97fadaa2fbd623c65badeefadc", "e7420b6492b8c876d23cd8a1156e35af4bc5dc5703fb4b79b376cb268a718e2e"],
      tags: ["eth", "integr"],
      companionNetworks: {amb: 'integr/amb'},
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
      accounts: ["0x80f702eb861f36fe8fbbe1a7ccceb04ef7ddef714604010501a5f67c8065d446", "0x5b18f0adcca221f65373b20158f95313ecd51bde42b96a4c16f5eb851576bc06"],
      tags: ["amb", "integr"],
      hardfork: "byzantium",
    },
  },

  namedAccounts: {
    owner: 0,
    proxyAdmin: 1,
    admin: 2,
    relay: 3,
    bridge: 4,
    user: 5,
  },

  verify: {
    etherscan: {
      apiKey: "DY4Z86MQ2D9E24C6HB98PTA79EKJ5TQIFX",
    },
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
      "Bridge$",
      ":ValidatorSet$",
      "sAMB",
      "IWrapper",
    ]
  }
};

export default config;
