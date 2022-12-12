import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {addNewTokensToBridge, options, parseNet, readConfig_} from "./utils/utils";
import {getAddresses} from "./utils/prod_addresses";

const BRIDGE_NAME = "BSC_BscBridge";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const isMainNet = parseNet(hre.network).stage === 'main'

  let configFile = readConfig_(hre.network);
  const tokenPairs = configFile.getTokenPairs("bsc", "amb")

  const ambNet = hre.companionNetworks['amb']
  const ambBridge = await ambNet.deployments.get('BSC_AmbBridge');

  const deployOptions = await options(hre, BRIDGE_NAME, tokenPairs,
    {
      sideBridgeAddress: ambBridge.address,
      wrappingTokenAddress: configFile.tokens.WBNB.addresses.bsc,
      timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
      lockTime: isMainNet ? 60 * 10 : 60,
      minSafetyBlocks: 10,
    }, []
  );

  // upgrade to untrustless-mpc; set watchdogs and fee_provider roles; remove DEFAULT_ADMIN_ROLE from deployer; remove RELAY_ROLE from old relay
  let {owner} = await hre.getNamedAccounts();
  const prod_addresses = getAddresses(BRIDGE_NAME);
  // @ts-ignore
  deployOptions.proxy.execute.onUpgrade = {
    methodName: "upgrade",
    args: [
      prod_addresses.watchdogsAddresses, // grand WATCHDOG_ROLEs
      isMainNet ? prod_addresses.feeProviderAddress: owner,  // grand FEE_PROVIDER_ROLE
      isMainNet ? prod_addresses.relayAddress : owner,  // grand RELAY_ROLE to new mpc relay

      isMainNet ? prod_addresses.adminAddress : owner, // remove DEFAULT_ADMIN_ROLE from this address
      isMainNet ? prod_addresses.feeProviderAddress: owner // remove RELAY_ROLE from this address (it's old relay address)
    ]
  };


  if (isMainNet) {
    console.log("To update prod contract remove this if statement :)");
  }  else {

    const deployResult = await hre.deployments.deploy(BRIDGE_NAME, {
      contract: BRIDGE_NAME,
      ...deployOptions,
    });

    configFile.bridges.bsc.side = deployResult.address;
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
func.tags = ["bridges_bsc"];
func.skip = async (hre: HardhatRuntimeEnvironment) => !hre.network.tags["bsc"]; // only bsc
