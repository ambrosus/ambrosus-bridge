import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {bridgesInNet, configPath, networkName, readConfig} from "./utils";
import { ethers } from "ethers";

const BRIDGE_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("BRIDGE_ROLE"));


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  let configFile = readConfig(configPath(hre.network));
  const bridgesInThisNetwork = bridgesInNet("eth", configFile)
  const netName = networkName(hre.network)

  const {owner} = await hre.getNamedAccounts();


  // set sideBridge to ambrosus bridge

  if (netName == "amb") {
    const ethBridge = await hre.companionNetworks['eth'].deployments.getOrNull('EthBridge');
    if (!ethBridge) throw new Error("[Setting sideBridgeAddress] Deploy EthBridge first")

    console.log("get sideBridge")
    const curAddr = await hre.deployments.read("AmbBridge", {from: owner}, 'sideBridgeAddress');
    if (curAddr != ethBridge.address)
      console.log("set sideBridge", curAddr, ethBridge.address)
      await hre.deployments.execute("AmbBridge",
        {from: owner, log: true},
        'setSideBridge', ethBridge.address
      );
  }

  // add bridge address to tokens
  console.log("add bridge address to tokens")

  for (const token of Object.values(configFile.tokens)) {
    if (!token.addresses[netName]) continue;  // not deployed
    if (token.primaryNet == netName) continue;  // it's not bridgeErc20, no need to set role

    const notSetBridges = await Promise.all(bridgesInThisNetwork.filter(async (br) => {
      return !await hre.deployments.read(
        token.symbol, {from: owner},
        "hasRole", BRIDGE_ROLE, br)
    }))

    if (notSetBridges.length > 0)
      await hre.deployments.execute(token.symbol, {from: owner, log: true},
        "setBridgeAddressesRole", notSetBridges)
  }


  // todo add new tokens to bridges
};


export default func;
func.tags = ["after_deploy"];
