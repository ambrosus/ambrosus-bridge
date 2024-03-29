import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {addNewTokensToBridge, options, parseNet, readConfig_} from "./utils/utils";

const BRIDGE_NAME = "ETH_EthBridge";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const isMainNet = parseNet(hre.network).stage === 'main'

  let configFile = readConfig_(hre.network);
  const tokenPairs = configFile.getTokenPairs("eth", "amb")

  const ambNet = hre.companionNetworks['amb']
  const ambBridge = await ambNet.deployments.get('ETH_AmbBridge');

  const deployOptions: any = await options(hre, BRIDGE_NAME, tokenPairs,
    {
      sideBridgeAddress: ambBridge.address,
      wrappingTokenAddress: configFile.tokens.WETH.networks.eth.address,
      timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
      lockTime: isMainNet ? 60 * 10 : 60,
      minSafetyBlocks: isMainNet ? 10 : 2,
    },
    []
  )

  if (isMainNet) {
    console.log("To update prod contract remove this if statement :)");
  } else {

    const deployResult = await hre.deployments.deploy(BRIDGE_NAME, {
      contract: BRIDGE_NAME,
      ...deployOptions
    });

    configFile.bridges.eth.side = deployResult.address;
    configFile.save()

    if (deployResult.newlyDeployed) {
      console.log('Call this cmd second time to update tokens')
      return;
    }
  }

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, BRIDGE_NAME);
};


export default func;
func.tags = ["bridges_eth"];
func.skip = async (hre: HardhatRuntimeEnvironment) => !hre.network.tags["eth"]; // only eth
