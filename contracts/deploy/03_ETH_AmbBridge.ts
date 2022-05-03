import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {
  addNewTokensToBridge,
  networkType, options,
  readConfig, setSideBridgeAddress,
} from "./utils";


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  let configFile = readConfig(hre.network);

  const tokenPairs = configFile.getTokenPairs("amb", "eth")

  const deployResult = await hre.deployments.deploy("ETH_AmbBridge", {
    contract: "ETH_AmbBridge",
    ...await options(hre, tokenPairs,
      {
        sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
        wrappingTokenAddress: configFile.tokens.SAMB.addresses.amb,
        fee: 1000,  // todo
        timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
        lockTime: isMainNet ? 60 * 10 : 60,
        minSafetyBlocks: 10,
      },
      [
        isMainNet ? 13_000_000_000 : 0  // minimum difficulty
      ]
    )
  });

  configFile.bridges.eth.amb = deployResult.address;
  configFile.save()

  if (deployResult.newlyDeployed) {
    console.log('Call this cmd second time to set sideBridgeAddress or update tokens')
    return;
  }


  // set sideBridgeAddress
  await setSideBridgeAddress("ETH_AmbBridge", configFile.bridges.eth.side, hre)

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, "ETH_AmbBridge");
};

export default func;
func.tags = ["bridges_eth"];
