import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {
  addNewTokensToBridge,
  configPath,
  getTokenPairs,
  networkType,
  readConfig,
  writeConfig
} from "./utils";


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  const path = configPath(hre.network);
  let configFile = readConfig(path);

  const {owner, proxyAdmin} = await hre.getNamedAccounts();
  // todo get admin and relay from getNamedAccounts
  const admin = owner;
  const relay = owner;

  const tokenPairs = getTokenPairs("amb", "eth", hre.network)

  const deployResult = await hre.deployments.deploy("ETH_AmbBridge", {
    contract: "ETH_AmbBridge",
    from: owner,
    proxy: {
      owner: proxyAdmin,
      proxyContract: "proxyTransparent",
      execute: {
        init: {
          methodName: "initialize",
          args: [
            {
              sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
              adminAddress: admin,
              relayAddress: relay,
              wrappingTokenAddress: configFile.tokens.SAMB.addresses.amb,
              tokenThisAddresses: Object.keys(tokenPairs),
              tokenSideAddresses: Object.values(tokenPairs),
              fee: 1000,  // todo
              feeRecipient: owner,   // todo
              timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
              lockTime: isMainNet ? 60 * 10 : 60,
              minSafetyBlocks: 10,
            },
            isMainNet ? 13_000_000_000 : 0  // minimum difficulty
          ]
        }
      }
    },
    log: true,
  });

  configFile.bridges.eth.amb = deployResult.address;
  writeConfig(path, configFile);

  if (deployResult.newlyDeployed) {
    console.log('Call this cmd second time to set sideBridgeAddress or update tokens')
    return;
  }


  // set sideBridgeAddress
  const ethBridge = configFile.bridges.eth.side;
  if (!ethBridge) {
    console.log("[Setting sideBridgeAddress] Deploy ETH_EthBridge first")
    return
  }

  const curAddr = await hre.deployments.read("ETH_AmbBridge", {from: owner}, 'sideBridgeAddress');
  if (curAddr != ethBridge) {
    console.log("[Setting sideBridgeAddress] old", curAddr, "new", ethBridge)
    await hre.deployments.execute("ETH_AmbBridge", {from: owner, log: true}, 'setSideBridge', ethBridge);
  }

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, "ETH_AmbBridge");
};

export default func;
func.tags = ["bridges_eth"];
