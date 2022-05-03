import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {networkName, readConfig} from "./utils";


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const netName = networkName(hre.network);
  let configFile = readConfig(hre.network);
  const bridgesInThisNetwork = configFile.bridgesInNet(netName)

  const {owner} = await hre.getNamedAccounts();

  for (const token of Object.values(configFile.tokens)) {
    if (token.addresses[netName]) continue;  // already deployed
    if (token.primaryNet == netName) continue;  // it's not bridgeErc20

    const {address} = await hre.deployments.deploy(token.symbol, {
      contract: "BridgeERC20Test", // todo remove `Test` on prod
      from: owner,
      args: [token.name, token.symbol, token.denomination, bridgesInThisNetwork],
      log: true,
    });

    token.addresses[netName] = address;
    configFile.save();
  }


};



export default func;
func.tags = ["tokens"];
