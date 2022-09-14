// Will set BRIDGE_ROLE to newly deployed bridges
// on already deployed BridgeERC20 tokens


import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {getBridgesDecimals, parseNet, readConfig_} from "./utils/utils";
import {ethers} from "ethers";
import {isAddress} from "ethers/lib/utils";
import {isTokenNotBridgeERC20, Token} from "./utils/config";


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const netName = parseNet(hre.network).name
  let configFile = readConfig_(hre.network);
  const {owner} = await hre.getNamedAccounts();

  // add bridge address to tokens
  console.log("add bridge address to tokens")


  // can have more than 1 bridge addresses; can convert decimals between networks
  const addAmb = async (token: Token) => {
    const {bridgesAddresses, bridgesDecimals} = getBridgesDecimals(configFile, token);

    const bridgesAddressesToSet = [];
    const bridgesDecimalsToSet = [];

    for (let i = 0; i < bridgesAddresses.length; i++) {
      const realBridgeDecimals = await hre.deployments.read(token.symbol, {from: owner}, "sideTokenDecimals", bridgesAddresses[i])
      if (realBridgeDecimals != bridgesDecimals[i]) {
        bridgesAddressesToSet.push(bridgesAddresses[i]);
        bridgesDecimalsToSet.push(bridgesDecimals[i]);
      }
    }
    if (bridgesAddressesToSet.length > 0)
      await hre.deployments.execute(token.symbol, {from: owner, log: true},
        "setSideTokenDecimals", bridgesAddressesToSet, bridgesDecimalsToSet);
  }

  // more lightweight contract; only 1 bridge address
  const addNonAmb = async (token: Token) => {
    const bridgeAddress = configFile.bridges[netName].side || ethers.constants.AddressZero
    const realBridgeAddress = await hre.deployments.read(token.symbol, {from: owner}, "bridgeAddress");

    if (realBridgeAddress != bridgeAddress)
      await hre.deployments.execute(token.symbol, {from: owner, log: true},
        "setBridgeAddress", bridgeAddress);
  }


  for (const token of Object.values(configFile.tokens)) {
    if (!isAddress(token.addresses[netName])) continue;  // not deployed
    if (isTokenNotBridgeERC20(token, netName)) continue;  // it's not bridgeErc20, no need to set role

    if (hre.network.tags["amb"]) await addAmb(token);
    else await addNonAmb(token);
  }

};


export default func;
func.tags = ["tokens_add_bridges"];
