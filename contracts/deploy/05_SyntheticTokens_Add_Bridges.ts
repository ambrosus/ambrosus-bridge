// Will set BRIDGE_ROLE to newly deployed bridges
// on already deployed BridgeERC20 tokens


import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {getBridgesDecimals, parseNet, readConfig_} from "./utils/utils";
import {ethers} from "ethers";
import {isAddress} from "ethers/lib/utils";
import {Config, isTokenPrimary, Token} from "./utils/config";


const BRIDGE_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("BRIDGE_ROLE"));


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const netName = parseNet(hre.network).name
  let configFile = readConfig_(hre.network);
  const {owner} = await hre.getNamedAccounts();

  // add bridge address to tokens
  console.log("add bridge address to tokens")


  // try to use old BridgeERC20 method for backward compatibility with legacy contracts
  const addLegacy = async (token: Token) => {
    const notSetBridges = (await Promise.all(
      bridgesInNet(netName, configFile)
        .map(async (br) => {
          const hasRole = await hre.deployments.read(token.symbol, {from: owner}, "hasRole", BRIDGE_ROLE, br)
          return hasRole ? null : br
        })))
      .filter(v => v != null)


    if (notSetBridges.length > 0)
      await hre.deployments.execute(token.symbol, {from: owner, log: true},
        "setBridgeAddressesRole", notSetBridges)
  }

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
    if (!token.isActive) continue;
    if (!isAddress(token.networks[netName]?.address)) continue;  // not deployed
    if (isTokenPrimary(token, netName)) continue;  // it's not synthetic token, no need to set role

    try {
      await addLegacy(token) // try to use old method first
    } catch (e) {
      if (netName == "amb") {
        await addAmb(token)
      } else {
        await addNonAmb(token)
      }
    }

  }

};

// get all deployed bridges in `net` network;
// for amb it's array of amb addresses for each network pair (such "amb-eth" or "amb-bsc")
// for other networks is array of one address
function bridgesInNet(net: string, configFile: Config): string[] {
  const bridges = (net == "amb") ?
    Object.values(configFile.bridges).map(i => i.amb) :
    [configFile.bridges[net].side];
  return bridges.filter(i => !!i);  // filter out empty strings
}


export default func;
func.tags = ["tokens_add_bridges"];
