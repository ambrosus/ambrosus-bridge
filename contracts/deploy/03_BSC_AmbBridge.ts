import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {
  addNewTokensToBridge,
  getBscValidators,
  options,
  parseNet,
  readConfig_,
  setSideBridgeAddress
} from "./utils/utils";

const BRIDGE_NAME = "BSC_AmbBridge";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const isMainNet = parseNet(hre.network).stage === 'main'

  let configFile = readConfig_(hre.network);
  const tokenPairs = configFile.getTokenPairs("amb", "bsc")

  const bscNet = hre.companionNetworks['bsc'];
  const [initialEpoch, initialValidators] = await getBscValidators(bscNet);
  const chainId = await bscNet.getChainId();

  if (isMainNet) {
    console.log("To update prod contract remove this if statement :)");

  }  else {

    const deployResult = await hre.deployments.deploy(BRIDGE_NAME, {
      contract: BRIDGE_NAME,
      ...await options(hre, BRIDGE_NAME, tokenPairs,
        {
          sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
          wrappingTokenAddress: configFile.tokens.SAMB.addresses.amb,
          timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
          lockTime: isMainNet ? 60 * 10 : 60,
          minSafetyBlocks: 10,
        },
        [
          initialValidators,
          initialEpoch,
          parseInt(chainId),
        ],
      )
    });

    configFile.bridges.bsc.amb = deployResult.address;
    configFile.save()

    if (deployResult.newlyDeployed) {
      console.log('Call this cmd second time to set sideBridgeAddress or update tokens')
      return;
    }

  }

  // set sideBridgeAddress
  await setSideBridgeAddress(BRIDGE_NAME, configFile.bridges.bsc.side, hre)

  // add new tokens

  // DISABLE WBNB TOKEN ( will be along with enabling USDC coin for gas economy :) )
  // todo remove this after call
  if (parseNet(hre.network).stage === "main") {
    tokenPairs["0xA96C522fA8Df99BB73A6E317A1afb0E3FA13b735"] = "0x0000000000000000000000000000000000000000"
  }
  console.log(tokenPairs);
  // END

  await addNewTokensToBridge(tokenPairs, hre, BRIDGE_NAME);
};


export default func;
func.tags = ["bridges_bsc"];
func.skip = async (hre: HardhatRuntimeEnvironment) => !hre.network.tags["amb"]; // only amb
