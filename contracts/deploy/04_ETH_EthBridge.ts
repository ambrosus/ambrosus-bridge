import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "ethers";
import {addNewTokensToBridge, getAmbValidators, options, parseNet, readConfig_} from "./utils/utils";
import {getAddresses} from "./utils/prod_addresses";

const BRIDGE_NAME = "ETH_EthBridge";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const isMainNet = parseNet(hre.network).stage === 'main'

  let configFile = readConfig_(hre.network);
  const tokenPairs = configFile.getTokenPairs("eth", "amb")

  const ambNet = hre.companionNetworks['amb']
  const ambBridge = await ambNet.deployments.get('ETH_AmbBridge');

  let {owner} = await hre.getNamedAccounts();
  const optionsWithOnUpgrade: any = await options(hre, BRIDGE_NAME, tokenPairs,
    {
      sideBridgeAddress: ambBridge.address,
      wrappingTokenAddress: ethers.constants.AddressZero, // todo WETH address
      timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
      lockTime: isMainNet ? 60 * 10 : 60,
      minSafetyBlocks: isMainNet ? 10 : 2,
    },
    []
  )

  // Uncomment when upgrade is needed
  // const prod_addresses = getAddresses(BRIDGE_NAME);
  // optionsWithOnUpgrade.proxy.execute.onUpgrade = {
  //   methodName: "upgrade",
  //   args: [
  //     [ // watchdogs (user-relays + master-relay + admin)
  //       "0x260cfE305cA40CaE1a32Ba7611137eF4d7146233", // Kevin
  //       "0xEB1c6a8a84063B1cef8B9a23AB87Bf926035A21a", // Lang
  //       "0x40B7d71E70fA6311cB0b300c1Ba6926A2A9000b8", // Rory
  //       "0xb017DcCC473499C83f1b553bE564f3CeAf002254", // Andrey
  //
  //       isMainNet ? prod_addresses.relayAddress : owner, // Master relay
  //       isMainNet ? prod_addresses.adminAddress : owner, // Admin
  //     ],
  //     isMainNet ? prod_addresses.relayAddress : owner, // fee_provider (relay)
  //     isMainNet ? prod_addresses.adminAddress : owner, // remove DEFAULT_ADMIN_ROLE from this address
  //   ]
  // };

  if (isMainNet) {
    console.log("To update prod contract remove this if statement :)");

  } else {

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
  }

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, BRIDGE_NAME);
};


export default func;
func.tags = ["bridges_eth"];
func.skip = async (hre: HardhatRuntimeEnvironment) => !hre.network.tags["eth"]; // only eth
