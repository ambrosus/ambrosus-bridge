import path from "path";
import fs from "fs";
import {HardhatRuntimeEnvironment} from "hardhat/types";


interface Token {
  name: string;
  symbol: string;
  denomination: number;
  addresses: { [net: string]: string }
  primaryNet: string;
}

interface Config {
  tokens: { [symb: string]: Token };
  bridges: { [net: string]: { amb: string, side: string } };
}


export function networkName(network: any): string {
  for (const networkName of ['amb', 'eth'])
    if (network.tags[networkName])
      return networkName;
  throw "Network missing networkName tag";
}

export function networkType(network: any): string {
  for (const networkType of ['testnet', 'mainnet'])
    if (network.tags[networkType])
      return networkType;
  throw "Network missing networkType tag";
}


export function getTokenPairs(thisNet: string, sideNet: string, network: any): { [k: string]: string } {
  return _getTokensPair(thisNet, sideNet, readConfig(configPath(network)));
}

function _getTokensPair(thisNet: string, sideNet: string, configFile: Config): { [k: string]: string } {
  const tokensPair: { [k: string]: string } = {};

  for (const tokenThis of Object.values(configFile.tokens)) {
    if (!tokenThis.addresses[thisNet] || !tokenThis.addresses[sideNet]) continue;
    tokensPair[tokenThis.addresses[thisNet]] = tokenThis.addresses[sideNet];
  }
  return tokensPair;
}


export async function addNewTokensToBridge(tokenPairs: { [k: string]: string },
                                           hre: HardhatRuntimeEnvironment,
                                           bridgeName: string): Promise<void> {
  const {owner} = await hre.getNamedAccounts();

  const newTokens = Object.entries(tokenPairs)
    .filter(async (pair) => {
      const tokenSide = await hre.deployments.read(bridgeName, {from: owner}, 'tokenAddresses', pair[0]);
      return tokenSide != pair[1]
    })
  if (newTokens.length == 0) return;

  const newTokens_ = Object.fromEntries(newTokens);
  await hre.deployments.execute(bridgeName, {from: owner, log: true},
    'tokensAdd', Object.keys(newTokens_), Object.values(newTokens_)
  )

}


// get all deployed bridges in `net` network;
// for amb it's array of amb addresses for each network pair (such "amb-eth" or "amb-bsc")
// for other networks is array of one address
export function bridgesInNet(net: string, configFile: Config): string[] {
  const bridges = (net == "amb") ?
    Object.values(configFile.bridges).map(i => i.amb) :
    [configFile.bridges[net].side];
  return bridges.filter(i => !!i);  // filter out empty strings
}

export function configPath(network: any): string {
  return path.resolve(__dirname, `../config-${networkType(network)}.json`);
}

export function writeConfig(path: string, config: Config) {
  fs.writeFileSync(path, JSON.stringify(config, null, 2));
}

export function readConfig(tokenPath: string): Config {
  return require(tokenPath);
}


// :(((
export function urlFromHHProvider(provider: any): string {
  while (provider && !provider.url) provider = provider._wrapped;
  return provider.url
}

