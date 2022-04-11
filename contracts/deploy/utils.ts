import path from "path";
import fs from "fs";


interface Tokens {
  name: string;
  symbol: string;
  address: string;
  addressesOnOtherNetworks: {[net: string]: string}
}

interface Config {
  tokens: { [net:string]: Tokens[] };
  bridges: { [net:string]: {amb: string, side: string} };
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


export function getTokensPair(thisNet: string, sideNet: string, network: any): any {
  return _getTokensPair(thisNet, sideNet, readConfig(configPath(network)));
}

function _getTokensPair(thisNet: string, sideNet: string, configFile: Config): any {
  const tokensPair: { [k: string]: string } = {};

  for (const tokenThis of configFile.tokens[thisNet]) {
    if (!tokenThis.address || !tokenThis.addressesOnOtherNetworks[sideNet]) continue;
    tokensPair[tokenThis.address] = tokenThis.addressesOnOtherNetworks[sideNet];
  }

  for (const tokenSide of configFile.tokens[sideNet]) {
    if (!tokenSide.address || !tokenSide.addressesOnOtherNetworks[thisNet]) continue;
    tokensPair[tokenSide.addressesOnOtherNetworks[thisNet]] = tokenSide.address;
  }

  return [Object.keys(tokensPair), Object.values(tokensPair)];
}

export function bridgesInNet(net: string, configFile: Config): string[] {
  if (net == "amb")
    return Object.values(configFile.bridges).map(i => i.amb)

  return [configFile.bridges[net].side];
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

