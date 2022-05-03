import {EthereumProvider, HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {
  addNewTokensToBridge,
  configPath,
  getTokenPairs,
  networkType,
  readConfig, setSideBridgeAddress, urlFromHHProvider,
  writeConfig
} from "./utils";
import vsAbi from "../abi/ValidatorSet.json";


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  const path = configPath(hre.network);
  let configFile = readConfig(path);

  const {owner, proxyAdmin} = await hre.getNamedAccounts();
  // todo get admin and relay from getNamedAccounts
  const admin = owner;
  const relay = owner;

  const tokenPairs = getTokenPairs("amb", "bsc", hre.network)

  const bscNet = hre.companionNetworks['bsc']
  const [initialEpoch, initialValidators] = await getValidators(bscNet.provider);
  const chainId = bscNet.getChainId();

  const deployResult = await hre.deployments.deploy("BSC_AmbBridge", {
    contract: "BSC_AmbBridge",
    from: owner,
    proxy: {
      owner: proxyAdmin,
      proxyContract: "proxyTransparent",
      execute: {
        init: {
          methodName: "initialize",
          args: [
            {
              sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
              adminAddress: admin,
              relayAddress: relay,
              wrappingTokenAddress: configFile.tokens.SAMB.addresses.amb,
              tokenThisAddresses: Object.keys(tokenPairs),
              tokenSideAddresses: Object.values(tokenPairs),
              fee: 1000,  // todo
              feeRecipient: owner,   // todo
              timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
              lockTime: isMainNet ? 60 * 10 : 60,
              minSafetyBlocks: 10,
            },
            initialEpoch,
            initialValidators,
            chainId,
          ]
        }
      }
    },
    log: true,
  });

  configFile.bridges.eth.amb = deployResult.address;
  writeConfig(path, configFile);

  if (deployResult.newlyDeployed) {
    console.log('Call this cmd second time to set sideBridgeAddress or update tokens')
    return;
  }


  // set sideBridgeAddress
  await setSideBridgeAddress("BSC_AmbBridge", configFile.bridges.bsc.side, hre)

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, "BSC_AmbBridge");
};




async function getValidators(bscProvider: EthereumProvider): Promise<[number, string[]]> {
  const provider = new ethers.providers.JsonRpcProvider(urlFromHHProvider(bscProvider))
  const {number: block} = await provider.getBlock('latest');
  const epoch = block / 200;
  const epochStart = epoch * 200;
  const blockWithValidators = await provider.getBlock(epochStart);

  // todo get validators from blockWithValidators
  const validators = [""];

  return [epoch, validators];
}



export default func;
func.tags = ["bridges_bnb"];
