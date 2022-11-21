import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {addNewTokensToBridge, options, parseNet, readConfig_, setSideBridgeAddress} from "./utils/utils";

const BRIDGE_NAME = "ETH_AmbBridge";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const isMainNet = parseNet(hre.network).stage === 'main'

  let configFile = readConfig_(hre.network);
  const tokenPairs = configFile.getTokenPairs("amb", "eth")

  const deployOptions: any = await options(hre, BRIDGE_NAME, tokenPairs,
    {
      sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
      wrappingTokenAddress: configFile.tokens.SAMB.networks.amb.address,
      timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
      lockTime: isMainNet ? 60 * 10 : 60,
      minSafetyBlocks: 10,
    },
    [ // untrusless1 params:
      5, // threshold
      [  // multisigners
        "0x260cfE305cA40CaE1a32Ba7611137eF4d7146233", // Kevin
        "0xEB1c6a8a84063B1cef8B9a23AB87Bf926035A21a", // Lang
        "0x40B7d71E70fA6311cB0b300c1Ba6926A2A9000b8", // Rory
        "0xb017DcCC473499C83f1b553bE564f3CeAf002254", // Andrey
        "0x0f071e1785e3E115360E04c9C8D53e958E6f85FE", // Master
      ]
    ],
  )


  if (isMainNet) {
    console.log("To update prod contract remove this if statement :)");
  }  else {

    const deployResult = await hre.deployments.deploy(BRIDGE_NAME, {
      contract: BRIDGE_NAME,
      ...deployOptions
    });

    configFile.bridges.eth.amb = deployResult.address;
    configFile.save()

    if (deployResult.newlyDeployed) {
      console.log('Call this cmd second time to set sideBridgeAddress or update tokens')
      return;
    }

  }

  // set sideBridgeAddress
  await setSideBridgeAddress(BRIDGE_NAME, configFile.bridges.eth.side, hre)

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, BRIDGE_NAME);
};

export default func;
func.tags = ["bridges_eth"];
func.skip = async (hre: HardhatRuntimeEnvironment) => !hre.network.tags["amb"]; // only amb
