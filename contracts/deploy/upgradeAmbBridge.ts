import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {
    configPath,
    getTokenPairs,
    networkType,
    readConfig,
} from "./utils";


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    if (hre.network.live && !hre.network.tags["amb"]) return;
    const isMainNet = networkType(hre.network) === 'mainnet'

    const path = configPath(hre.network);
    let configFile = readConfig(path);

    const {owner, proxyAdmin} = await hre.getNamedAccounts();
    // todo get admin and relay from getNamedAccounts
    const admin = owner;
    const relay = owner;

    const tokenPairs = getTokenPairs("amb", "eth", hre.network)

    await hre.deployments.deploy("AmbBridge", {
        from: proxyAdmin,
        proxy: {
            proxyContract: "proxyTransparent",
            execute: {
                init: {
                    methodName: "initialize",
                    args: [{
                        sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
                        adminAddress: admin,
                        relayAddress: relay,
                        wrappingTokenAddress: configFile.tokens.SAMB.addresses.amb,
                        tokenThisAddresses: Object.keys(tokenPairs),
                        tokenSideAddresses: Object.values(tokenPairs),
                        fee: 1000,  // todo
                        feeRecipient: owner,   // todo
                        timeframeSeconds: isMainNet ? 14400 : 1,
                        lockTime: isMainNet ? 1000 : 1,
                        minSafetyBlocks: 10,
                    }]
                }
            }
        },
        log: true,
    });
};

export default func;
func.tags = ["upgradeambbridge"];
