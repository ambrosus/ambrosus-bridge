import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {networkName, readConfig, configPath, writeConfig, bridgesInNet} from "./utils";


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const path = configPath(hre.network);
  let configFile = readConfig(path);

  let tokens = configFile.tokens[networkName(hre.network)];

  const bridgesInThisNetwork = bridgesInNet("eth", configFile)

  const {owner} = await hre.getNamedAccounts();

  for (const token of tokens) {
    if (token.address) continue;  // already deployed

    const {address} = await hre.deployments.deploy(token.name, {
      contract: "BridgeERC20",
      from: owner,
      args: [token.name, token.symbol, bridgesInThisNetwork],
      log: true,
    });

    token.address = address;
    writeConfig(path, configFile);
  }


};



export default func;
func.tags = ["BridgeERC20Tokens"];
