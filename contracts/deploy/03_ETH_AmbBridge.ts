import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {addNewTokensToBridge, networkType, options, readConfig, setSideBridgeAddress} from "./utils";

const BRIDGE_NAME = "ETH_AmbBridge";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  let configFile = readConfig(hre.network);
  const tokenPairs = configFile.getTokenPairs("amb", "eth")

  const deployResult = await hre.deployments.deploy(BRIDGE_NAME, {
    contract: BRIDGE_NAME,
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
      ],
      [
        ["0xd17C51395268d1EE25bDc1F2E107b7E11b41f5d1", "0x295C2707319ad4BecA6b5bb4086617fD6F240CfE"],
        2
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
  await setSideBridgeAddress(BRIDGE_NAME, configFile.bridges.eth.side, hre)

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, BRIDGE_NAME);
};

export default func;
func.tags = ["bridges_eth"];
