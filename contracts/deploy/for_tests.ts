import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.name !== "hardhat") return;
  const {owner} = await hre.getNamedAccounts();

  const commonArgs = {
    sideBridgeAddress: ethers.constants.AddressZero,
    relayAddress: ethers.constants.AddressZero,
    tokenThisAddresses: [
      "0x0000000000000000000000000000000000000001",
      "0x0000000000000000000000000000000000000003",
      "0x0000000000000000000000000000000000000005",
    ],
    tokenSideAddresses: [
      "0x0000000000000000000000000000000000000002",
      "0x0000000000000000000000000000000000000004",
      "0x0000000000000000000000000000000000000006",
    ],
    fee: 1000,
    feeRecipient: ethers.constants.AddressZero,
    timeframeSeconds: 14400,
    lockTime: 1000,
    minSafetyBlocks: 10,
  };

  await hre.deployments.deploy("CommonBridge", {
    from: owner,
    args: [commonArgs],
    log: true,
  });

  await hre.deployments.deploy("AmbBridgeTest", {
    from: owner,
    args: [commonArgs, ethers.constants.AddressZero],
    log: true,
  });

  await hre.deployments.deploy("EthBridgeTest", {
    contract: "EthBridgeTest",
    from: owner,
    args: [commonArgs,
      [ // initial validators
        "0x0000000000000000000000000000000000000001",
        "0x0000000000000000000000000000000000000002",
      ],
      ethers.constants.AddressZero, // validatorSetAddress_
      ethers.constants.HashZero, // lastProcessedBlock
    ],
    log: true,
  });

  await hre.deployments.deploy("MockERC20", {
    contract: "MockERC20",
    from: owner,
    args: [
      [ethers.constants.AddressZero], // bridgeAddresses
    ],
    log: true,
  });

  await hre.deployments.deploy("wAMB", {
    contract: "wAMB",
    from: owner,
    args: ["wAMB", "wAmb"],
    log: true,
  });
};

export default func;

func.tags = ["for_tests"];
