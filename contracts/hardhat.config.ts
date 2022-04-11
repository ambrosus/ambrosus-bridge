import "@nomiclabs/hardhat-waffle";
import "hardhat-deploy";
import "hardhat-gas-reporter";
import "solidity-coverage";
import "hardhat-abi-exporter";


import * as dotenv from "dotenv";

dotenv.config();
const PK = [process.env.PRIVATEKEY || "00000000"];

module.exports = {
    networks: {
        hardhat: {
            hardfork: "byzantium",
            companionNetworks: {
                amb: 'hardhat',
                eth: 'hardhat',
            },
        },
        eth_test: {
            url: "https://ropsten.infura.io/v3/" + process.env.INFURA_KEY,
            accounts: PK,
            tags: ["eth", "testnet"],
            companionNetworks: {
                amb: 'amb_test',
            },
        },
        eth_main: {
            url: "https://mainnet.infura.io/v3/" + process.env.INFURA_KEY,
            accounts: PK,
            tags: ["eth", "mainnet"],
            companionNetworks: {
                amb: 'amb_main',
            },
        },
        amb_test: {
            url: "https://network.ambrosus-test.io",
            accounts: PK,
            tags: ["amb", "testnet"],
            hardfork: "byzantium",
            companionNetworks: {
                eth: 'eth_test',
            },
        },
        amb_main: {
            url: "https://network.ambrosus.io",
            accounts: PK,
            tags: ["amb", "mainnet"],
            hardfork: "byzantium",
            companionNetworks: {
                eth: 'eth_main',
            },
        }
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
                        runs: 200,
                    },
                    // Note: for amb deploy
                    evmVersion: "byzantium"
                },
            },{
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
        ]
    }
};
