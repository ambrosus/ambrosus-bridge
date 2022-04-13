import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {addNewTokensToBridge, configPath, getTokenPairs, networkType, readConfig, writeConfig} from "./utils";

const relayAddress = "0x295c2707319ad4beca6b5bb4086617fd6f240cfe" // todo get from something?

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  const path = configPath(hre.network);
  let configFile = readConfig(path);

  const {owner} = await hre.getNamedAccounts();
  const tokenPairs = getTokenPairs("amb", "eth", hre.network)

  const deployResult = await hre.deployments.deploy("AmbBridge", {
    from: owner,
    args: [
      {
        sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
        relayAddress: relayAddress,
        tokenThisAddresses: Object.keys(tokenPairs),
        tokenSideAddresses: Object.values(tokenPairs),
        fee: 1000,  // todo
        feeRecipient: owner,   // todo
        timeframeSeconds: isMainNet ? 14400 : 1,
        lockTime: isMainNet ? 1000 : 1,
        minSafetyBlocks: 10,
      },
      configFile.tokens.wAMB.addresses["amb"],
    ],
    log: true,
    skipIfAlreadyDeployed: true
  });

  configFile.bridges.eth.amb = deployResult.address;
  writeConfig(path, configFile);


  if (deployResult.newlyDeployed) {
    console.log('Call this cmd second time to set sideBridgeAddress or update tokens')
    return;
  }


  // set sideBridgeAddress
  const ethBridge = configFile.bridges.eth.side;
  if (!ethBridge) throw new Error("[Setting sideBridgeAddress] Deploy EthBridge first")

  const curAddr = await hre.deployments.read("AmbBridge", {from: owner}, 'sideBridgeAddress');
  if (curAddr != ethBridge)
    await hre.deployments.execute("AmbBridge", {from: owner, log: true}, 'setSideBridge', ethBridge);


  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, "AmbBridge");

};

export default func;
func.tags = ["bridges"];
