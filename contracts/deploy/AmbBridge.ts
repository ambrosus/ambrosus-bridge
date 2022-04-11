import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";

const relayAddress = "0x295c2707319ad4beca6b5bb4086617fd6f240cfe" // todo get from something?

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;

  const {owner} = await hre.getNamedAccounts();
  const {address: ambWrapper} = await hre.deployments.get('wAMB');


  const t = hre.deployments.deploy("AmbBridge", {
    from: owner,
    args: [
      {
        sideBridgeAddress: null, // amb deployed before eth
        relayAddress: relayAddress,
        tokenThisAddresses: [],
        tokenSideAddresses: [],
        fee: 1000,  // todo
        feeRecipient: owner,   // todo
        timeframeSeconds: hre.network.live ? 14400 : 1,
        lockTime: hre.network.live ? 1000 : 1,
        minSafetyBlocks: 10,
      },
      ambWrapper
    ],
    log: true,
  });
  console.log(t)
};

export default func;
func.tags = ["AmbBridge"];
func.dependencies = ['wAMB'];
