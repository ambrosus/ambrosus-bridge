import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {getTokensPair} from "./utils";

const relayAddress = "0x295c2707319ad4beca6b5bb4086617fd6f240cfe" // todo get from something?

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;

  const {owner} = await hre.getNamedAccounts();
  const {address: ambWrapper} = await hre.deployments.get('wAMB');

  const [tokensThis, tokensSide] = getTokensPair("amb", "eth", hre.network)

  const deployResult = await hre.deployments.deploy("AmbBridge", {
    from: owner,
    args: [
      {
        sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
        relayAddress: relayAddress,
        tokenThisAddresses: tokensThis,
        tokenSideAddresses: tokensSide,
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


  if (deployResult.newlyDeployed) {
    console.log("Run ambBridge deploy after ethBridge deploy to set sideBridgeAddress")
  } else {
    const ethBridge = await hre.companionNetworks['eth'].deployments.getOrNull('EthBridge');
    if (!ethBridge) throw new Error("[Setting sideBridgeAddress] Deploy EthBridge first")

    const setAddr = await hre.deployments.read("AmbBridge", {from: owner}, 'sideBridgeAddress');
    if (setAddr == ethBridge.address) return;

    await hre.deployments.execute("AmbBridge",
      {from: owner, log: true},
      'setSideBridge', ethBridge.address
    );
    console.log("ambBridge.sideBridgeAddress set to", ethBridge.address)
  }

};

export default func;
func.tags = ["AmbBridge"];
func.dependencies = ['wAMB'];
