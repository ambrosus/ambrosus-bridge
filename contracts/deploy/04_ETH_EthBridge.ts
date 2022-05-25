import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "ethers";
import {addNewTokensToBridge, getAmbValidators, networkType, options, readConfig} from "./utils";

const BRIDGE_NAME = "ETH_EthBridge";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["eth"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  let configFile = readConfig(hre.network);
  const tokenPairs = configFile.getTokenPairs("eth", "amb")

  const ambNet = hre.companionNetworks['amb']
  const {address: sideBridgeAddress} = await ambNet.deployments.get('ETH_AmbBridge');

  const optionsWithOnUpgrade: any = {
    ...await options(
         hre,
         tokenPairs,
         {
           sideBridgeAddress: sideBridgeAddress,
           wrappingTokenAddress: configFile.tokens.WETH?.addresses.eth || ethers.constants.AddressZero,
           transferFeeRecipient: ethers.constants.AddressZero,
           bridgeFeeRecipient: ethers.constants.AddressZero,
           timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
           lockTime: isMainNet ? 60 * 10 : 60,
           minSafetyBlocks: isMainNet ? 10 : 2,
         },
        await getAmbValidators(ambNet.provider, hre),
    )
  }

  // Uncomment when updateLastProcessedBlock is needed
  // optionsWithOnUpgrade.proxy.execute.onUpgrade = {
  //   methodName: "updateLastProcessedBlock",
  //   args: ["0x1111111111111111111111111111111111111111111111111111111111111111"]
  // };


  const deployResult = await hre.deployments.deploy(BRIDGE_NAME, {
    contract: BRIDGE_NAME,
    ...optionsWithOnUpgrade
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
func.tags = ["bridges_eth"];
