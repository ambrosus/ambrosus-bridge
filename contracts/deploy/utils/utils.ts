import {HardhatRuntimeEnvironment, Network} from "hardhat/types";
import {DeployOptions} from "hardhat-deploy/types";
import {Config, readConfig, Token} from "./config";
import {getAddresses} from "./prod_addresses";
import {isAddress} from "ethers/lib/utils";

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
  for (const tokenThis of Object.keys(tokenPairs)) {
    const tokenSide = await hre.deployments.read(bridgeName, {from: owner}, 'tokenAddresses', tokenThis);
    if (tokenPairs[tokenThis] == tokenSide)
      delete tokenPairs[tokenThis];
  }

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

export async function options(hre: HardhatRuntimeEnvironment, bridgeName: string, tokenPairs: { [k: string]: string },
                              commonArgs: any, args: any[]): Promise<DeployOptions> {

  const network = parseNet(hre.network);
  let {owner} = await hre.getNamedAccounts();

  // on testnets use only 1 account for all roles;
  // multisig threshold == 1, so no upgrade confirmations needed
  const cfg = (network.stage === "main") ? getAddresses(bridgeName) :
    {
      adminAddress: owner,
      relayAddress: owner,
      feeProviderAddress: owner,
      watchdogsAddresses: [owner],
      transferFeeRecipient: owner,
      bridgeFeeRecipient: owner,
      multisig: {
        admins: [owner],
        threshold: 1
      }
    };

  if (owner != cfg.adminAddress) {
    throw `Deploying from address '${owner}', but config adminAddress is '${cfg.adminAddress}'`;
  }

  // add this args to user args
  const reallyCommonArgs = {
    relayAddress: cfg.relayAddress,
    feeProviderAddress: cfg.feeProviderAddress,
    watchdogsAddresses: cfg.watchdogsAddresses,
    transferFeeRecipient: cfg.transferFeeRecipient,
    bridgeFeeRecipient: cfg.bridgeFeeRecipient,

    tokenThisAddresses: Object.keys(tokenPairs),
    tokenSideAddresses: Object.values(tokenPairs),
  }
  // commonArgs is contract `ConstructorArgs` struct
  commonArgs = {...reallyCommonArgs, ...commonArgs};

  return {
    from: owner,
    proxy: {
      owner: owner,
      proxyArgs: ["{implementation}", "{data}", cfg.multisig.admins, cfg.multisig.threshold],
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


// get bridges and decimals for BridgeERC20_Amb contract
export function getBridgesDecimals(configFile: Config, token: Token) {
  const bridgesAddresses = [];
  const bridgesDecimals = [];

  for (const netName of Object.keys(configFile.bridges)) {
    if (token.networks[netName] === undefined) continue; // skip if token not deployed on this network

    const ambBridgeAddress = configFile.bridges[netName].amb;
    if (!isAddress(ambBridgeAddress)) continue; // skip non existing bridges :)

    const sideNetDecimals = token.networks[netName].denomination;
    bridgesAddresses.push(ambBridgeAddress);
    bridgesDecimals.push(sideNetDecimals);
  }

  return {bridgesAddresses, bridgesDecimals};
}
