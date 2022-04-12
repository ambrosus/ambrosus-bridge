import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {networkName, readConfig, configPath, writeConfig, bridgesInNet} from "./utils";


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const path = configPath(hre.network);
  let configFile = readConfig(path);

  const netName = networkName(hre.network);
  const bridgesInThisNetwork = bridgesInNet('eth', configFile)

  const {owner} = await hre.getNamedAccounts();

  for (const token of Object.values(configFile.tokens)) {
    if (token.addresses[netName]) continue;  // already deployed
    if (token.primaryNet == netName) continue;  // it's not bridgeErc20

    const {address} = await hre.deployments.deploy(token.symbol, {
      contract: "BridgeERC20Test", // todo remove `Test` on prod
      from: owner,
      args: [token.name, token.symbol, bridgesInThisNetwork],
      log: true,
    });

    token.addresses[netName] = address;
    writeConfig(path, configFile);
  }


};



export default func;
func.tags = ["tokens"];
