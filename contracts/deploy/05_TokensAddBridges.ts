import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {bridgesInNet, configPath, networkName, readConfig} from "./utils";
import { ethers } from "ethers";

const BRIDGE_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("BRIDGE_ROLE"));


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  let configFile = readConfig(configPath(hre.network));
  const netName = networkName(hre.network)
  const bridgesInThisNetwork = bridgesInNet(netName, configFile)

  const {owner} = await hre.getNamedAccounts();

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

};


export default func;
func.tags = ["tokens_add_bridges"];
