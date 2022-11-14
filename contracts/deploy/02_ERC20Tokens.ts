// Will deploy all tokens from config, except:
// - ambrosus network
// - already deployed tokens (that have address set on this network)
// - native tokens (primaryNet field == network.name)

import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {getBridgesDecimals, parseNet, readConfig_} from "./utils/utils";
import {isTokenNotBridgeERC20, Token} from "./utils/config";
import {ethers} from "hardhat";


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const netName = parseNet(hre.network).name;
  let configFile = readConfig_(hre.network);
  const {owner} = await hre.getNamedAccounts();

  // can have more than 1 bridge addresses; can convert decimals between networks
  const deployAmb = async (token: Token) => {
    const {bridgesAddresses, bridgesDecimals} = getBridgesDecimals(configFile, token);

    const {address} = await hre.deployments.deploy(token.symbol, {
      contract: "BridgeERC20_Amb",
      args: [token.name, token.symbol, token.denomination, bridgesAddresses, bridgesDecimals],
      from: owner, log: true,
    });
    return address;
  }

  // more lightweight contract; only 1 bridge address
  const deployNonAmb = async (token: Token) => {
    const bridgeAddress = configFile.bridges[netName].side || ethers.constants.AddressZero

    const {address} = await hre.deployments.deploy(token.symbol, {
      contract: "BridgeERC20",
      args: [token.name, token.symbol, token.denomination, bridgeAddress],
      from: owner, log: true,
    });
    return address;
  }


  for (const token of Object.values(configFile.tokens)) {
    if (!token.isActive) continue;
    if (token.addresses[netName] != "DEPLOY") continue;  // already deployed or shouldn't be deployed
    if (isTokenNotBridgeERC20(token, netName)) continue;  // it's not bridgeErc20

    const address = (hre.network.tags["amb"]) ?
      await deployAmb(token) :
      await deployNonAmb(token);

    token.addresses[netName] = address;
    configFile.save();
  }


};


export default func;
func.tags = ["tokens"];
