import "@nomiclabs/hardhat-waffle";
import "hardhat-deploy";
import "hardhat-gas-reporter";
import "solidity-coverage";
import "hardhat-abi-exporter";
import * as dotenv from "dotenv";
import {HardhatUserConfig, HttpNetworkUserConfig} from "hardhat/types";
import {ethers} from "ethers";


dotenv.config();
// 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, used instead of empty PK
const devPK = "34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5"

const bscScanApiKey = "NFH875QU828E37MQD7XB3QHFBE4XTC2AKH"
const polygonScanApiKey = "5B48B7UTG14J3UNYQGQ9I3UHIHB7PND7VE"
const optimismScanApiKey = "9RU98DN9AAUB99HZXMFN5WETSP68CEN75S"
const arbitrumScanApiKey = "9RU98DN9AAUB99HZXMFN5WETSP68CEN75S" // TODO: change api key

const sideNets = ["eth", "bsc", "polygon", "optimism", "arbitrum"];
const bscExtraFields = {verify: {etherscan: {apiKey: bscScanApiKey}}};
const polygonExtraFields = {verify: {etherscan: {apiKey: polygonScanApiKey}}};
const optimismExtraFields = {verify: {etherscan: {apiKey: optimismScanApiKey}}};
const arbitrumExtraFields = {verify: {etherscan: {apiKey: optimismScanApiKey}}};


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

    ...ambNetwork("dev", "https://network.ambrosus-dev.io"),
    ...ambNetwork("test", "https://network.ambrosus-test.io"),
    ...ambNetwork("main", "https://network.ambrosus.io"),

    ...network("eth", "dev", "https://goerli.infura.io/v3/" + process.env.INFURA_KEY),
    ...network("eth", "test", "https://sepolia.infura.io/v3/" + process.env.INFURA_KEY),
    ...network("eth", "main", "https://mainnet.infura.io/v3/" + process.env.INFURA_KEY),

    ...network("polygon", "test", "https://polygon-mumbai.g.alchemy.com/v2/" + process.env.ALCHEMY_KEY, polygonExtraFields),

    ...network("optimism", "test", "https://opt-goerli.g.alchemy.com/v2/6S14oqFC3IHZl0F1trsy_vN41gmDqDaF", optimismExtraFields),

    ...network("arbitrum", "test", "https://arb-goerli.g.alchemy.com/v2/ZOwlu6LylAMGi7hAN_LU5VD8fLqUnuH_", arbitrumExtraFields),

    ...network("bsc", "dev", "https://data-seed-prebsc-1-s1.binance.org:8545", bscExtraFields),
    ...network("bsc", "test", "https://data-seed-prebsc-1-s1.binance.org:8545", bscExtraFields),
    ...network("bsc", "main", "https://bsc-dataseed1.binance.org", bscExtraFields),

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
      ":ValidatorSet$",
      "sAMB",
      "IWrapper",
      ":ERC20$",
      "ProxyMultiSig$"
    ]
  }
};

function network(network: string, stage: string, url: string, extraFields: any = {}): { [name: string]: HttpNetworkUserConfig } {
  console.assert(["dev", "test", "main"].includes(stage), "stage must be dev/test/main");

  const account = (stage === "main") ?
    process.env["PRIVATEKEY_OWNER_" + network.toUpperCase()] :
    devPK;

  return {
    [`${stage}/${network}`]: {
      url: url,
      tags: [network],
      companionNetworks: {amb: `${stage}/amb`},
      accounts: [account || ethers.constants.HashZero],
      ...extraFields
    },
  }
}

function ambNetwork(stage: string, url: string, otherFields: any = {}): { [name: string]: HttpNetworkUserConfig } {
  const companionNetworks: { [net: string]: string } = {};
  sideNets.forEach(sideNet => companionNetworks[sideNet] = `${stage}/${sideNet}`);

  otherFields.hardfork = "byzantium";
  otherFields.companionNetworks = companionNetworks;

  return network("amb", stage, url, otherFields);
}

export default config;
