import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.name !== "hardhat") return;
  const {owner, admin} = await hre.getNamedAccounts();

  const {address: mockAddr} = await hre.deployments.deploy("BridgeERC20Test", {
    contract: "BridgeERC20Test",
    from: owner,
    args: [
      "Mock", "Mock", 18,
      [ethers.constants.AddressZero], // bridgeAddresses
    ],
    log: true,
  });

  const {address: wrapperAddr} = await hre.deployments.deploy("sAMB", {
    contract: "sAMB",
    from: owner,
    args: ["sAMB", "sAMB"],
    log: true,
  });


  const commonArgs = {
    sideBridgeAddress: ethers.constants.AddressZero,
    adminAddress: ethers.constants.AddressZero,
    relayAddress: ethers.constants.AddressZero,
    wrappingTokenAddress: wrapperAddr,
    tokenThisAddresses: [],
    tokenSideAddresses: [],
    transferFeeRecipient: ethers.constants.AddressZero,
    bridgeFeeRecipient: ethers.constants.AddressZero,
    timeframeSeconds: 14400,
    lockTime: 1000,
    minSafetyBlocks: 10,
  };

  await hre.deployments.deploy("CommonBridgeTest", {
    from: owner,
    args: [commonArgs],
    log: true,
  });

    await hre.deployments.deploy("CheckPoWTest", {
        from: owner,
        args: [
            0 // minimum difficulty
        ],
    });

    await hre.deployments.deploy("CheckAuraTest", {
        from: owner,
        args: [
            "0x0000000000000000000000000000000000000F00", // validatorSetAddress_
        ],
    });

  await hre.deployments.deploy("CheckPoSATest", {
    from: owner,
    args: [
      "0x61"  // chainId
    ],
    log: true,
  });

  await hre.deployments.deploy("ProxyMultisigTest", {
    from: owner,
    args: []
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
