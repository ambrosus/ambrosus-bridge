// Will set BRIDGE_ROLE to newly deployed bridges
// on already deployed BridgeERC20 tokens


import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {parseNet, readConfig_} from "./utils/utils";
import {ethers} from "ethers";
import { isAddress } from "ethers/lib/utils";

const BRIDGE_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("BRIDGE_ROLE"));


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const netName = parseNet(hre.network).name
  let configFile = readConfig_(hre.network);

  // list of deployed bridge (or bridgeS on amb) to which the BRIDGE_ROLE should be assigned
  const bridgesInThisNetwork = configFile.bridgesInNet(netName)

  const {owner} = await hre.getNamedAccounts();

  // add bridge address to tokens
  console.log("add bridge address to tokens")

  for (const token of Object.values(configFile.tokens)) {
    if (!isAddress(token.addresses[netName])) continue;  // not deployed
    if (token.primaryNet == netName) continue;  // it's not bridgeErc20, no need to set role

    const notSetBridges = (await Promise.all(
      bridgesInThisNetwork
        .map(async (br) => {
          const hasRole = await hre.deployments.read(token.symbol, {from: owner}, "hasRole", BRIDGE_ROLE, br)
          return hasRole ? null : br
        })))
      .filter(v => v != null)


    if (notSetBridges.length > 0)
      await hre.deployments.execute(token.symbol, {from: owner, log: true},
        "setBridgeAddressesRole", notSetBridges)
  }

};


export default func;
func.tags = ["tokens_add_bridges"];
