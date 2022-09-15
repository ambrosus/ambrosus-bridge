import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {getSideNetDecimalsOrTokenDenomination, parseNet, readConfig_} from "./utils/utils";
import {isAddress} from "ethers/lib/utils";
import {isTokenNotBridgeERC20} from "./utils/config";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const netName = parseNet(hre.network).name;
  const {owner} = await hre.getNamedAccounts();

  let configFile = readConfig_(hre.network);


  for (const token of Object.values(configFile.tokens)) {
    if (token.addresses[netName] != "DEPLOY") continue;  // already deployed or shouldn't be deployed
    if (!isTokenNotBridgeERC20(token, netName)) continue;  // it's bridgeErc20

    console.log("Deploying test mintable token", token.symbol, "on", netName);
    const {address} = await hre.deployments.deploy(token.symbol, {
      contract: "MintableERC20",
      from: owner,
      args: [token.name, token.symbol, getSideNetDecimalsOrTokenDenomination(token, netName)],
      log: true,
    });


    token.addresses[netName] = address;
    configFile.save();
  }

};

export default func;
func.tags = ["tokens"];
func.skip = async (hre: HardhatRuntimeEnvironment) => hre.network.live; // only testnets
