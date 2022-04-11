import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {networkName, readTokens, tokensPath, writeTokens} from "./utils";


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const tokenPath = tokensPath(hre.network);
  let tokensFile = readTokens(tokenPath);
  let tokens = tokensFile[networkName(hre.network)].tokens;

  const {owner} = await hre.getNamedAccounts();

  const bridgesInThisNetwork: string[] = [] // todo

  for (const token of tokens) {
    if (token.address) continue;  // already deployed

    const {address} = await hre.deployments.deploy(token.name, {
      contract: "BridgeERC20",
      from: owner,
      args: [token.name, token.symbol, bridgesInThisNetwork],
      log: true,
    });

    token.address = address;
    writeTokens(tokenPath, tokensFile);
  }


};



export default func;
func.tags = ["BridgeERC20Tokens"];
