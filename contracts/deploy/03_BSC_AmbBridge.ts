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
  // const [initialEpoch, initialValidators] = await getBscValidators(bscNet);
  const chainId = await bscNet.getChainId();
  const initialEpoch = 100500
  const initialValidators = ["0x049153b8dae0a232ac90d20c78f1a5d1de7b7dc5",
    "0x1284214b9b9c85549ab3d2b972df0deef66ac2c9",
    "0x35552c16704d214347f29fa77f77da6d75d7c752",
    "0x980a75ecd1309ea12fa2ed87a8744fbfc9b863d5",
    "0xa2959d3f95eae5dc7d70144ce1b73b403b7eb6e0",
    "0xb71b214cb885500844365e95cd9942c7276e7fd8",
    "0xf474cf03cceff28abc65c9cbae594f725c80e12d",
  ]

  const deployResult = await hre.deployments.deploy(BRIDGE_NAME, {
    contract: BRIDGE_NAME,
    ...await options(hre, tokenPairs,
      {
        sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
        wrappingTokenAddress: configFile.tokens.SAMB.addresses.amb,
        transferFeeRecipient: ethers.constants.AddressZero,
        bridgeFeeRecipient: ethers.constants.AddressZero,
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


  // set sideBridgeAddress
  await setSideBridgeAddress(BRIDGE_NAME, configFile.bridges.bsc.side, hre)

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, BRIDGE_NAME);
};


export default func;
func.tags = ["bridges_bsc"];
func.skip = async (hre: HardhatRuntimeEnvironment) => !hre.network.tags["amb"]; // only amb
