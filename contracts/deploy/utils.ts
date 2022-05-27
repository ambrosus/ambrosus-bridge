import path from "path";
import fs from "fs";
import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployOptions} from "hardhat-deploy/types";
import {ethers} from "ethers";
import vsAbi from "../abi/ValidatorSet.json";
import {Receipt} from "hardhat-deploy/dist/types";


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

  save(): void;

  getTokenPairs(thisNet: string, sideNet: string): { [k: string]: string }

  bridgesInNet(net: string): string[]
}


export function readConfig(network: any): Config {
  const tokenPath = path.resolve(__dirname, `../config-${networkType(network)}.json`)
  const config = require(tokenPath);

  config.save = () => fs.writeFileSync(tokenPath, JSON.stringify(config, null, 2));
  config.getTokenPairs = (thisNet: string, sideNet: string) => getTokenPairs(thisNet, sideNet, config)
  config.bridgesInNet = (net: string) => bridgesInNet(net, config)

  return config;
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


function getTokenPairs(thisNet: string, sideNet: string, configFile: Config): { [k: string]: string } {
  const tokenPair: { [k: string]: string } = {};

  for (const token of Object.values(configFile.tokens)) {

    if (token.addresses[thisNet] && token.addresses[sideNet])
      tokenPair[token.addresses[thisNet]] = token.addresses[sideNet];

    if (token.primaryNet === sideNet && token.nativeAnalog)   // native token for sideNet
      tokenPair[ethers.constants.AddressZero] = token.addresses[thisNet];

  }

  return tokenPair;
}


// get all deployed bridges in `net` network;
// for amb it's array of amb addresses for each network pair (such "amb-eth" or "amb-bsc")
// for other networks is array of one address
function bridgesInNet(net: string, configFile: Config): string[] {
  const bridges = (net == "amb") ?
    Object.values(configFile.bridges).map(i => i.amb) :
    [configFile.bridges[net].side];
  return bridges.filter(i => !!i);  // filter out empty strings
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


export async function options(hre: HardhatRuntimeEnvironment, tokenPairs: { [k: string]: string },
                              commonArgs: any, args: any[]): Promise<DeployOptions> {

  const {owner, proxyAdmin} = await hre.getNamedAccounts();
  // todo get admin and relay from getNamedAccounts
  const admin = owner;
  const relay = owner;

  // add this args to user args
  const reallyCommonArgs = {
    adminAddress: admin,
    relayAddress: relay,
    feeRecipient: owner,   // todo
    tokenThisAddresses: Object.keys(tokenPairs),
    tokenSideAddresses: Object.values(tokenPairs),
  }
  // commonArgs is contract `ConstructorArgs` struct
  commonArgs = {...reallyCommonArgs, ...commonArgs};

  return {
    from: owner,
    proxy: {
      owner: owner,
      proxyArgs: ["{implementation}", "{data}", [owner, proxyAdmin], 2],
      proxyContract: "ProxyMultiSig",
      execute: {
        init: {
          methodName: "initialize",
          args: [commonArgs, ...args]
        }
      }
    },
    log: true
  }
}

async function getValidators(network: any, vsAddress: string, vsAbi: any, blockNum: number,): Promise<string[]> {
  const provider = new ethers.providers.JsonRpcProvider(urlFromHHProvider(network.provider));
  const vsContract = ethers.ContractFactory.getContract(vsAddress, vsAbi);
  return await vsContract.connect(provider).getValidators({blockTag: blockNum});
}

export async function getBscValidators(bscNetwork: any, bscDeployReceipt: Receipt): Promise<[number, string[]]> {
  const vsAddress = "0x0000000000000000000000000000000000001000";
  const validators = await getValidators(bscNetwork, vsAddress, vsAbi, bscDeployReceipt.blockNumber);
  const epoch = bscDeployReceipt.blockNumber / 200;

  return [epoch, validators];
}

export async function getAmbValidators(ambNetwork: any, ambDeployReceipt: Receipt): Promise<[string[], string, string]> {
  const vsAddress = "0x0000000000000000000000000000000000000F00";
  const validators = await getValidators(ambNetwork, vsAddress, vsAbi, ambDeployReceipt.blockNumber);

  return [validators, vsAddress, ambDeployReceipt.blockHash]
}


// :(((
export function urlFromHHProvider(provider: any): string {
  while (provider && !provider.url) provider = provider._wrapped;
  return provider.url
}
