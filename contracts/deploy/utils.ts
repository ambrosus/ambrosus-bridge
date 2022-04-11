import path from "path";
import fs from "fs";

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
  return _getTokensPair(thisNet, sideNet, readTokens(tokensPath(network)));
}

function _getTokensPair(thisNet: string, sideNet: string, tokensFile: any): any {
  const tokensPair: { [k: string]: string } = {};

  for (const tokenThis of tokensFile[thisNet].tokens) {
    if (!tokenThis.address) continue;
    if (!tokenThis.addressesOnOtherNetworks[sideNet]) continue;
    tokensPair[tokenThis.address] = tokenThis.addressesOnOtherNetworks[sideNet];
  }

  for (const tokenSide of tokensFile[sideNet].tokens) {
    if (!tokenSide.address) continue;
    if (!tokenSide.addressesOnOtherNetworks[thisNet]) continue;
    tokensPair[tokenSide.addressesOnOtherNetworks[thisNet]] = tokenSide.address;
  }

  return [Object.keys(tokensPair), Object.values(tokensPair)];
}

export function tokensPath(network: any): string {
  return path.resolve(__dirname, `../tokens-${networkType(network)}.json`);
}

export function writeTokens(tokenPath: string, tokens: any) {
  fs.writeFileSync(tokenPath, JSON.stringify(tokens, null, 2));
}

export function readTokens(tokenPath: string): any {
  return require(tokenPath);
}


// :(((
export function urlFromHHProvider(provider: any): string {
  while (provider && !provider.url) provider = provider._wrapped;
  return provider.url
}

