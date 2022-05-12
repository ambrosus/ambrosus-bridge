import {EthereumProvider, HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {
  addNewTokensToBridge,
  networkType,
  options,
  readConfig,
  setSideBridgeAddress,
  urlFromHHProvider,
} from "./utils";

const BRIDGE_NAME = "BSC_AmbBridge";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  let configFile = readConfig(hre.network);
  const tokenPairs = configFile.getTokenPairs("amb", "bsc")

  const bscNet = hre.companionNetworks['bsc']
  const [initialEpoch, initialValidators] = await getBscValidators(bscNet.provider);
  const chainId = await bscNet.getChainId();

  const deployResult = await hre.deployments.deploy(BRIDGE_NAME, {
    contract: BRIDGE_NAME,
    ...await options(hre, tokenPairs,
      {
        sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
        wrappingTokenAddress: configFile.tokens.SAMB.addresses.amb,
        fee: 1000,  // todo
        timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
        lockTime: isMainNet ? 60 * 10 : 60,
        minSafetyBlocks: 10,
      },
      [
        initialEpoch,
        initialValidators,
        chainId,
      ],
    )
  });

  configFile.bridges.eth.amb = deployResult.address;
  configFile.save()

  if (deployResult.newlyDeployed) {
    console.log('Call this cmd second time to set sideBridgeAddress or update tokens')
    return;
  }


  // set sideBridgeAddress
  await setSideBridgeAddress(BRIDGE_NAME, configFile.bridges.bsc.side, hre)

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, BRIDGE_NAME);
};


async function getBscValidators(bscProvider: EthereumProvider): Promise<[number, string[]]> {
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
