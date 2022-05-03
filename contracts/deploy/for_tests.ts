import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.name !== "hardhat") return;
  const {owner, proxyAdmin} = await hre.getNamedAccounts();

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
    fee: 1000,
    feeRecipient: ethers.constants.AddressZero,
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
    log: true,
  });

  await hre.deployments.deploy("CheckAuraTest", {
    from: owner,
    args: [
      [
        "0x4c9785451bb2CA3E91B350C06bcB5f974cA33F79",
        "0x90B2Ce3741188bCFCe25822113e93983ecacfcA0",
        "0xAccdb7a2268BC4Af0a1898e725138888ba1Ca6Fc"
      ],
      "0x0000000000000000000000000000000000000F00", // validatorSetAddress_
    ],
    log: true,
  });

  await hre.deployments.deploy("CheckPoSATest", {
    from: owner,
    args: [
      [
        "0x049153b8dae0a232ac90d20c78f1a5d1de7b7dc5",
        "0x35552c16704d214347f29fa77f77da6d75d7c752",
        "0x7a1a4ad9cc746a70ee58568466f7996dd0ace4e8",
        "0x980a75ecd1309ea12fa2ed87a8744fbfc9b863d5",
        "0xa2959d3f95eae5dc7d70144ce1b73b403b7eb6e0",
        "0xb71b214cb885500844365e95cd9942c7276e7fd8",
        "0xc89c669357d161d57b0b255c94ea96e179999919",
        "0xf474cf03cceff28abc65c9cbae594f725c80e12d"
      ],
      64643,
      "0x61"
    ],
    log: true,
  });
};

export default func;

func.tags = ["for_tests"];
