import path from "path";
import fs from "fs";
import {HardhatRuntimeEnvironment} from "hardhat/types";
import {ethers} from "ethers";


interface Token {
  name: string;
  symbol: string;
  denomination: number;
  addresses: { [net: string]: string }
  primaryNet: string;
  nativeAnalog: string | null;
}

interface Config {
  tokens: { [symb: string]: Token };
  bridges: { [net: string]: { amb: string, side: string } };
}


export function networkName(network: any): string {
  const r = ['amb', 'eth'].find(t => network.tags[t]);
  if (!r) throw "Network missing networkName tag";
  return r
}

export function networkType(network: any): string {
  const r = ['devnet', 'testnet', 'mainnet', 'integr'].find(t => network.tags[t]);
  if (!r) throw "Network missing networkType tag";
  return r
}


export function getTokenPairs(thisNet: string, sideNet: string, network: any): { [k: string]: string } {
  return _getTokensPair(thisNet, sideNet, readConfig(configPath(network)));
}

function _getTokensPair(thisNet: string, sideNet: string, configFile: Config): { [k: string]: string } {
  const tokensPair: { [k: string]: string } = {};

  for (const token of Object.values(configFile.tokens)) {

    if (token.addresses[thisNet] && token.addresses[sideNet])
      tokensPair[token.addresses[thisNet]] = token.addresses[sideNet];

    if (token.primaryNet === sideNet && token.nativeAnalog)   // native token for sideNet
      tokensPair[ethers.constants.AddressZero] = token.addresses[thisNet];

  }

  return tokensPair;
}


export async function addNewTokensToBridge(tokenPairs: { [k: string]: string },
                                           hre: HardhatRuntimeEnvironment,
                                           bridgeName: string): Promise<void> {
  const {owner} = await hre.getNamedAccounts();

  // remove from tokenPairs all tokens that are already in the bridge
  await Promise.all(Object.keys(tokenPairs).map(async (tokenThis) => {
    const tokenSide = await hre.deployments.read(bridgeName, {from: owner}, 'tokenAddresses', tokenThis);
    if (tokenPairs[tokenThis] == tokenSide)
      delete tokenPairs[tokenThis];
  }));

  if (Object.keys(tokenPairs).length == 0) {
    console.log("No new tokens to add to bridge");
    return;
  }

  console.log("Adding new tokens to bridge:", tokenPairs);
  await hre.deployments.execute(bridgeName, {from: owner, log: true},
    'tokensAddBatch', Object.keys(tokenPairs), Object.values(tokenPairs)
  )

}

export async function setSideBridgeAddress(deploymentName: string, sideAddress: string, hre: HardhatRuntimeEnvironment) {
  if (!sideAddress) {
    console.log(`[Setting sideBridgeAddress] Deploy side bridge for ${deploymentName} first`)
    return
  }
  const {owner} = await hre.getNamedAccounts();

  const curAddr = await hre.deployments.read(deploymentName, {from: owner}, 'sideBridgeAddress');
  if (curAddr != sideAddress) {
    console.log("[Setting sideBridgeAddress] old", curAddr, "new", sideAddress)
    await hre.deployments.execute(deploymentName, {from: owner, log: true}, 'setSideBridge', sideAddress);
  }
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

