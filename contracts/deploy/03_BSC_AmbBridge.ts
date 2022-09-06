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


  // set sideBridgeAddress
  await setSideBridgeAddress(BRIDGE_NAME, configFile.bridges.bsc.side, hre)

  // add new tokens

  // BUGFIX
  // todo remove this after redeploy and register bugged tokens
  if (parseNet(hre.network).stage === "dev") {
    tokenPairs["0x2FA390265Fd43c996f6335293ad60446494c621B"] = "0x0000000000000000000000000000000000000000"
    tokenPairs["0x6C5750Be93De7a4A3072bA10A6610C23e8399df1"] = "0x0000000000000000000000000000000000000000"
  } else if (parseNet(hre.network).stage === "test") {
    tokenPairs["0x4798Cbd108e3D7b531ef2f3d67E7fFabdEe29867"] = "0x0000000000000000000000000000000000000000"
    tokenPairs["0x63B825C40a78e2e9A7aeaC83027215A022b37B93"] = "0x0000000000000000000000000000000000000000"
  } else if (parseNet(hre.network).stage === "main") {
    tokenPairs["0xA96C522fA8Df99BB73A6E317A1afb0E3FA13b735"] = "0x0000000000000000000000000000000000000000"
  }
  console.log(tokenPairs);
  // END OF BUGFIX

  await addNewTokensToBridge(tokenPairs, hre, BRIDGE_NAME);
};


export default func;
func.tags = ["bridges_bsc"];
func.skip = async (hre: HardhatRuntimeEnvironment) => !hre.network.tags["amb"]; // only amb
