import "@nomiclabs/hardhat-waffle";
import "hardhat-deploy";
import "hardhat-gas-reporter";
import "solidity-coverage";
import "hardhat-abi-exporter";
import * as dotenv from "dotenv";
import {HardhatUserConfig} from "hardhat/types";
import {ethers} from "ethers";


dotenv.config();
const devPK = process.env.PRIVATE_KEY_DEV || ethers.constants.HashZero;  // used for dev networks
const bscScanApiKey = process.env.ETHERSCAN_API_KEY;
const etherScanApiKey = process.env.ETHERSCAN_API_KEY;

const config: HardhatUserConfig = {

  networks: {
    hardhat: {
      blockGasLimit: 40000000, // amb value
      hardfork: "petersburg",
      companionNetworks: {amb: 'hardhat'},
      initialDate: "13 May 2022 18:10:36 GMT",

      // for debugging: run local node with state from some network
      forking: {
        enabled: false,
        url: "https://network.ambrosus-dev.io",
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
      accounts: [process.env.PRIVATEKEY_OWNER_ETH || ethers.constants.HashZero],
    },
    "main/eth": {
      url: "https://mainnet.infura.io/v3/" + process.env.INFURA_KEY,
      tags: ["eth", "mainnet"],
      companionNetworks: {amb: 'main/amb'},
      accounts: [process.env.PRIVATEKEY_OWNER_ETH || ethers.constants.HashZero],
    },
    "integr/eth": {
      url: "http://127.0.0.1:8502",
      accounts: [process.env.PRIVATEKEY_INTEGR1 || ethers.constants.HashZero, process.env.PRIVATEKEY_INTEGR2 || ethers.constants.HashZero],
      tags: ["eth", "integr"],
      companionNetworks: {amb: 'integr/amb'},
    },

    "dev/amb": {
      url: "https://network.ambrosus-dev.io",
      tags: ["amb", "devnet"],
      hardfork: "petersburg",
      companionNetworks: {eth: 'dev/eth', bsc: 'dev/bsc'},
      accounts: [devPK], // todo devPk
    },
    "test/amb": {
      url: "https://network.ambrosus-test.io",
      tags: ["amb", "testnet"],
      hardfork: "petersburg",
      companionNetworks: {eth: 'test/eth', bsc: 'test/bsc'},
      accounts: [process.env.PRIVATEKEY_OWNER_AMB || ethers.constants.HashZero],
    },
    "main/amb": {
      url: "https://network.ambrosus.io",
      tags: ["amb", "mainnet"],
      hardfork: "petersburg",
      companionNetworks: {eth: 'main/eth', bsc: 'main/bsc'},
      accounts: [process.env.PRIVATEKEY_OWNER_AMB || ethers.constants.HashZero],
    },
    "integr/amb": {
      url: "http://127.0.0.1:8545",
      accounts: [process.env.PRIVATEKEY_INTEGR3 || ethers.constants.HashZero, process.env.PRIVATEKEY_INTEGR4 || ethers.constants.HashZero],
      tags: ["amb", "integr"],
      hardfork: "petersburg",
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
      apiKey: etherScanApiKey,
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
          evmVersion: "petersburg"
        },
      }, {
        version: "0.4.22",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200,
          },
          // Note: for amb deploy
          evmVersion: "petersburg"
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
