import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "ethers";
import {addNewTokensToBridge, getAmbValidators, options, parseNet, readConfig_} from "./utils/utils";

const BRIDGE_NAME = "BSC_BscBridge";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const isMainNet = parseNet(hre.network).stage === 'main'

  let configFile = readConfig_(hre.network);
  const tokenPairs = configFile.getTokenPairs("bsc", "amb")

  const ambNet = hre.companionNetworks['amb']
  const ambBridge = await ambNet.deployments.get('BSC_AmbBridge');

  const deployResult = await hre.deployments.deploy(BRIDGE_NAME, {
    contract: BRIDGE_NAME,
    ...await options(hre, tokenPairs,
      {
        sideBridgeAddress: ambBridge.address,
        wrappingTokenAddress: configFile.tokens.WBNB.addresses.bsc,
        transferFeeRecipient: ethers.constants.AddressZero,
        bridgeFeeRecipient: ethers.constants.AddressZero,
        timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
        lockTime: isMainNet ? 60 * 10 : 60,
        minSafetyBlocks: isMainNet ? 10 : 2,
      },
      await getAmbValidators(ambNet),
    )
  });


  configFile.bridges.bsc.side = deployResult.address;
  configFile.save()

  if (deployResult.newlyDeployed) {
    console.log('Call this cmd second time to update tokens')
    return;
  }

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, BRIDGE_NAME);
};


export default func;
func.tags = ["bridges_bsc"];
func.skip = async (hre: HardhatRuntimeEnvironment) => !hre.network.tags["bsc"]; // only bsc
