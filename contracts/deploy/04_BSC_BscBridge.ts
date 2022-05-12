import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "ethers";
import {addNewTokensToBridge, getAmbValidators, networkType, options, readConfig} from "./utils";

const BRIDGE_NAME = "BSC_BscBridge";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["bsc"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  let configFile = readConfig(hre.network);
  const tokenPairs = configFile.getTokenPairs("bsc", "amb")

  const ambNet = hre.companionNetworks['amb']
  const {address: sideBridgeAddress} = await ambNet.deployments.get('BSC_AmbBridge');

  const deployResult = await hre.deployments.deploy(BRIDGE_NAME, {
    contract: BRIDGE_NAME,
    ...await options(hre, tokenPairs,
      {
        sideBridgeAddress: sideBridgeAddress,
        wrappingTokenAddress: configFile.tokens.WBNB?.addresses.eth || ethers.constants.AddressZero,
        fee: 10,    // todo
        timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
        lockTime: isMainNet ? 60 * 10 : 60,
        minSafetyBlocks: isMainNet ? 10 : 2,
      },
      await getAmbValidators(ambNet.provider),
    )
  });


  configFile.bridges.eth.side = deployResult.address;
  configFile.save()

  if (deployResult.newlyDeployed) {
    console.log('Call this cmd second time to update tokens')
    return;
  }

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, BRIDGE_NAME);
};


export default func;
func.tags = ["bridges_bnb"];
