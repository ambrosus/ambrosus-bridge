
import {HardhatRuntimeEnvironment, Network} from "hardhat/types";
import {DeployOptions} from "hardhat-deploy/types";
import {ethers} from "ethers";
import vsAbi from "../../abi/ValidatorSet.json";
import {Block} from "@ethersproject/abstract-provider";
import {Config, readConfig} from "./config";

export function readConfig_(network: Network): Config {
  return readConfig(parseNet(network).stage);
}

export function parseNet(network: Network): { stage: string; name: string } {
  if (network.name == "hardhat")
    throw "Hardhat network not supported"
  const [stage, name] = network.name.split('/')
  return {stage, name};
}

// actions

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

//

export async function options(hre: HardhatRuntimeEnvironment, tokenPairs: { [k: string]: string },
                              commonArgs: any, args: any[]): Promise<DeployOptions> {

  let {owner, admin, relay} = await hre.getNamedAccounts();

  // multisig admins and threshold
  let msAdmins = process.env.MULTISIG_ADDRESSES!.split(',');
  let msThresh = +process.env.MULTISIG_THRESHOLD!;

  // on testnets use only 1 account for all roles;
  // multisig threshold == 1, so no upgrade confirmations needed
  if (parseNet(hre.network).stage != "main") {
    [admin, relay] = [owner, owner];
    [msAdmins, msThresh] = [[owner], 1];
  }

  // add this args to user args
  const reallyCommonArgs = {
    adminAddress: admin,
    relayAddress: relay,
    transferFeeRecipient: owner, // todo
    bridgeFeeRecipient: owner,  // todo
    tokenThisAddresses: Object.keys(tokenPairs),
    tokenSideAddresses: Object.values(tokenPairs),
  }
  // commonArgs is contract `ConstructorArgs` struct
  commonArgs = {...reallyCommonArgs, ...commonArgs};

  return {
    from: owner,
    proxy: {
      owner: owner,
      proxyArgs: ["{implementation}", "{data}", msAdmins, msThresh],
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

// valildators

async function getValidatorsAndLatestBlock(network: any, vsAddress: string, vsAbi: any): Promise<[string[], Block, any]> {
  const provider = new ethers.providers.JsonRpcProvider(urlFromHHProvider(network.provider));
  const vsContract = ethers.ContractFactory.getContract(vsAddress, vsAbi).connect(provider);

  const latestBlock = await provider.getBlock('latest');
  const validators = await vsContract.getValidators({blockTag: "latest"});
  return [validators, latestBlock, vsContract];
}

export async function getBscValidators(bscNetwork: any): Promise<[number, string[]]> {
  const vsAddress = "0x0000000000000000000000000000000000001000";
  const [validators, latestBlock] = await getValidatorsAndLatestBlock(bscNetwork, vsAddress, vsAbi);
  const epoch = Math.floor(latestBlock.number / 200);

  return [epoch, validators];
}

export async function getAmbValidators(ambNetwork: any): Promise<[string[], string, string]> {
  const vsAddress = "0x0000000000000000000000000000000000000F00";
  const [validators, latestBlock, vsContract] = await getValidatorsAndLatestBlock(ambNetwork, vsAddress, vsAbi);

  // check that current validators match with the latest finalized event
  const logs = await vsContract.queryFilter(vsContract.filters.InitiateChange())
  const latestLog = logs[logs.length-1]
  const latestSet = vsContract.interface.parseLog(latestLog).args.newSet
  console.assert(JSON.stringify(latestSet) == JSON.stringify(validators),
    `ValidatorSet extracted from ${latestBlock.number} block doesn't equal to 
    ValidatorSet emitted in ${latestLog.blockNumber} block. 
    Probably, latest event doesn't finalized yet and this can cause a trouble.
    Try again at block ~${latestLog.blockNumber + latestSet.length/2}`)

  return [validators, vsAddress, latestBlock.hash]
}


// :(((
export function urlFromHHProvider(provider: any): string {
  while (provider && !provider.url) provider = provider._wrapped;
  return provider.url
}
