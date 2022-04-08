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
            hardfork: "byzantium"

        },
        rinkeby: {
            url: "https://rinkeby.infura.io/v3/" + process.env.INFURA_KEY,
            accounts: PK,
        },
        mainnet: {
            url: "https://mainnet.infura.io/v3/" + process.env.INFURA_KEY,
            accounts: PK,
        },
        amb: {
            url: "https://network.ambrosus-dev.io",
            accounts: PK,
            hardfork: "byzantium"
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
