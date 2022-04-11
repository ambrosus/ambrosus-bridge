import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import fs from "fs";
import path from "path";



const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const tokenPath = path.resolve(__dirname, `../tokens-${networkType(hre.network)}.json`);
  let tokensFile = readTokens(tokenPath);
  let tokens = tokensFile[networkName(hre.network)].tokens;

  const {owner} = await hre.getNamedAccounts();

  const bridgesInThisNetwork: string[] = [] // todo

  for (const token of tokens) {
    if (token.address) continue;  // already deployed

    const {address} = await hre.deployments.deploy(token.name, {
      from: owner,
      args: [token.name, token.symbol, bridgesInThisNetwork],
      log: true,
    });

    token.address = address;
    writeTokens(tokenPath, tokensFile);
  }


};


function networkName(network: any): string {
  for (const networkName of ['amb', 'eth'])
    if (network.tags[networkName])
      return networkName;
  throw "Network missing networkName tag";
}

function networkType(network: any): string {
  for (const networkType of ['testnet', 'mainnet'])
    if (network.tags[networkType])
      return networkType;
  throw "Network missing networkType tag";
}


function writeTokens(tokenPath: string, tokens: any) {
  fs.writeFileSync(tokenPath, JSON.stringify(tokens, null, 2));
}

function readTokens(tokenPath: string): any {
  return require(tokenPath);
}


export default func;
func.tags = ["BridgeERC20Tokens"];
