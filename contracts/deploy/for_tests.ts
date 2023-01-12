import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {owner, admin} = await hre.getNamedAccounts();

    const {address: mockAddr} = await hre.deployments.deploy("BridgeERC20Test", {
        from: owner,
        args: [
            "Mock", "Mock", 18,
            ethers.constants.AddressZero, // bridgeAddress
        ],
    });
    await hre.deployments.deploy("BridgeERC20_AmbTest", {
        from: owner,
        args: [
            "Mock", "Mock", 18,
            [ethers.constants.AddressZero], // bridgeAddresses
            [0], // bridgeDecimals
        ],
    });

    const {address: wrapperAddr} = await hre.deployments.deploy("sAMB", {
        from: owner,
        args: ["sAMB", "sAMB"],
    });


    await hre.deployments.deploy("CommonBridgeTest", {from: owner}); // can't use calldata in normal constructor, so ...
    await hre.deployments.execute("CommonBridgeTest", {from: owner}, "constructor_", {
        sideBridgeAddress: ethers.constants.AddressZero,
        relayAddress: ethers.constants.AddressZero,
        feeProviderAddress: ethers.constants.AddressZero,
        watchdogsAddresses: [ethers.constants.AddressZero],
        wrappingTokenAddress: wrapperAddr,
        tokenThisAddresses: [],
        tokenSideAddresses: [],
        transferFeeRecipient: ethers.constants.AddressZero,
        bridgeFeeRecipient: ethers.constants.AddressZero,
        timeframeSeconds: 14400,
        lockTime: 1000,
        minSafetyBlocks: 10,
    });

    await hre.deployments.deploy("CheckUntrustlessTest", {
        from: owner,
    });

    await hre.deployments.deploy("Faucet", {
        from: owner,
        args: [[owner]],
    });


    await hre.deployments.deploy("ProxyMultisigTest", {
        from: owner,
    });

    await hre.deployments.deploy("ProxyMultiSig", {
        from: owner,
        args: [
            mockAddr,
            "0x",
            [owner, admin],
            2
        ],
    });
};

export default func;
func.tags = ["for_tests"];
func.skip = async (hre: HardhatRuntimeEnvironment) => hre.network.name !== "hardhat"; // only hardhat
