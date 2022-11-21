import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {addNewTokensToBridge, options, parseNet, readConfig_, setSideBridgeAddress} from "./utils/utils";
import {getAddresses} from "./utils/prod_addresses";

const BRIDGE_NAME = "ETH_AmbBridge";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const isMainNet = parseNet(hre.network).stage === 'main'

  let configFile = readConfig_(hre.network);
  const tokenPairs = configFile.getTokenPairs("amb", "eth")

  const deployOptions: any = await options(hre, BRIDGE_NAME, tokenPairs,
    {
      sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
      wrappingTokenAddress: configFile.tokens.SAMB.addresses.amb,
      timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
      lockTime: isMainNet ? 60 * 10 : 60,
      minSafetyBlocks: 10,
    }, []
  )

  // upgrade to untrustless-mpc; set watchdogs and fee_provider roles; remove DEFAULT_ADMIN_ROLE from deployer
  let {owner} = await hre.getNamedAccounts();
  const prod_addresses = getAddresses(BRIDGE_NAME);
  deployOptions.proxy.execute.onUpgrade = {
    methodName: "upgrade",
    args: [
      [ // watchdogs (user-relays + master-relay + admin)
        "0x260cfE305cA40CaE1a32Ba7611137eF4d7146233", // Kevin
        "0xEB1c6a8a84063B1cef8B9a23AB87Bf926035A21a", // Lang
        "0x40B7d71E70fA6311cB0b300c1Ba6926A2A9000b8", // Rory
        "0xb017DcCC473499C83f1b553bE564f3CeAf002254", // Andrey

        isMainNet ? prod_addresses.relayAddress : owner, // Master relay
        isMainNet ? prod_addresses.adminAddress : owner, // Admin
      ],
      isMainNet ? prod_addresses.relayAddress : owner, // fee_provider (relay)
      isMainNet ? prod_addresses.adminAddress : owner, // remove DEFAULT_ADMIN_ROLE from this address
      'TODO' // new mpc relay address
      // todo remove RELAY_ROLE from old relay address
    ]
  };


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
