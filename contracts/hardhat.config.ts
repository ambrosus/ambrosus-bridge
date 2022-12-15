import "@nomiclabs/hardhat-waffle";
import "hardhat-deploy";
import "hardhat-gas-reporter";
import "solidity-coverage";
import "hardhat-abi-exporter";
import * as dotenv from "dotenv";
import {HardhatUserConfig} from "hardhat/types";
import {ethers} from "ethers";


dotenv.config();
// 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, used instead of empty PK
const devPK = "34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5"
const bscScanApiKey = "NFH875QU828E37MQD7XB3QHFBE4XTC2AKH"

const config: HardhatUserConfig = {

  networks: {
    hardhat: {
      blockGasLimit: 40000000, // amb value
      hardfork: "byzantium",
      companionNetworks: {amb: 'hardhat'},
      initialDate: "13 May 2022 18:10:36 GMT",

      // for debugging: run local node with state from some network
      forking: {
        enabled: false,
        url: "https://network.ambrosus-dev.io",
        // url: "https://eth-rinkeby.alchemyapi.io/v2/e1F5R9XuWDU2-zCtzaMDg4Ybb5SuoEDA"
        blockNumber: 0xb00ba,
      }
    },

    "dev/eth": {
      url: "https://goerli.infura.io/v3/" + process.env.INFURA_KEY,
      tags: ["eth", "devnet"],
      companionNetworks: {amb: 'dev/amb'},
      gasPrice: 9000000000,
      accounts: [devPK],
    },
    "test/eth": {
      url: "https://sepolia.infura.io/v3/" + process.env.INFURA_KEY,
      tags: ["eth", "testnet"],
      companionNetworks: {amb: 'test/amb'},
      accounts: [devPK],
    },
    "main/eth": {
      url: "https://mainnet.infura.io/v3/" + process.env.INFURA_KEY,
      tags: ["eth", "mainnet"],
      companionNetworks: {amb: 'main/amb'},
      accounts: [process.env.PRIVATEKEY_OWNER_ETH || ethers.constants.HashZero],
    },
    "integr/eth": {
      url: "http://127.0.0.1:8502",
      accounts: ["0x51d098d8aee092622149d8f3a79cc7b1ce36ff97fadaa2fbd623c65badeefadc", "e7420b6492b8c876d23cd8a1156e35af4bc5dc5703fb4b79b376cb268a718e2e"],
      tags: ["eth", "integr"],
      companionNetworks: {amb: 'integr/amb'},
    },

    "dev/amb": {
      url: "https://network.ambrosus-dev.io",
      tags: ["amb", "devnet"],
      hardfork: "byzantium",
      companionNetworks: {eth: 'dev/eth', bsc: 'dev/bsc'},
      accounts: [devPK], // todo devPk
    },
    "test/amb": {
      url: "https://network.ambrosus-test.io",
      tags: ["amb", "testnet"],
      hardfork: "byzantium",
      companionNetworks: {eth: 'test/eth', bsc: 'test/bsc'},
      accounts: [process.env.PRIVATEKEY_OWNER_AMB || ethers.constants.HashZero],
    },
    "main/amb": {
      url: "https://network.ambrosus.io",
      tags: ["amb", "mainnet"],
      hardfork: "byzantium",
      companionNetworks: {eth: 'main/eth', bsc: 'main/bsc'},
      accounts: [process.env.PRIVATEKEY_OWNER_AMB || ethers.constants.HashZero],
    },
    "integr/amb": {
      url: "http://127.0.0.1:8545",
      accounts: ["0x80f702eb861f36fe8fbbe1a7ccceb04ef7ddef714604010501a5f67c8065d446", "0x5b18f0adcca221f65373b20158f95313ecd51bde42b96a4c16f5eb851576bc06"],
      tags: ["amb", "integr"],
      hardfork: "byzantium",
    },


    "dev/bsc": {
      url: "https://data-seed-prebsc-1-s1.binance.org:8545/",
      tags: ["bsc", "devnet"],
      companionNetworks: {amb: 'dev/amb'},
      accounts: [devPK],
      verify: {
        etherscan: {
          apiKey: bscScanApiKey
        }
      },
    },
    "test/bsc": {
      url: "https://data-seed-prebsc-1-s1.binance.org:8545/",
      tags: ["bsc", "testnet"],
      companionNetworks: {amb: 'test/amb'},
      accounts: [process.env.PRIVATEKEY_OWNER_BSC || ethers.constants.HashZero], // todo devPk
      verify: {
        etherscan: {
          apiKey: bscScanApiKey
        }
      },
    },
    "main/bsc": {
      url: "https://bsc-dataseed1.binance.org/",
      tags: ["bsc", "mainnet"],
      companionNetworks: {amb: 'main/amb'},
      accounts: [process.env.PRIVATEKEY_OWNER_BSC || ethers.constants.HashZero],
      verify: {
        etherscan: {
          apiKey: bscScanApiKey
        }
      },
    },
  },

  namedAccounts: {
    owner: 0,
    admin: 1,
    relay: 2,
    bridge: 3,
    user: 4,
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
      "Bridge_",
      "sAMB",
      "IWrapper",
      ":ERC20$",
      "ProxyMultiSig$",
      "Faucet"
    ]
  }
};

export default config;
