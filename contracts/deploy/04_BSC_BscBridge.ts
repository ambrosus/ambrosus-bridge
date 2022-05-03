import {EthereumProvider, HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import vsAbi from "../abi/ValidatorSet.json";
import {ethers} from "ethers";
import {
  addNewTokensToBridge,
  networkType, options,
  readConfig,
  urlFromHHProvider,
} from "./utils";

const vsAddress = "0x0000000000000000000000000000000000000F00" // todo get from something?


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["bsc"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  let configFile = readConfig(hre.network);

  const tokenPairs = configFile.getTokenPairs("bsc", "amb")

  const ambNet = hre.companionNetworks['amb']
  const {address: sideBridgeAddress} = await ambNet.deployments.get('BSC_AmbBridge');
  const initialValidators = await getValidators(ambNet.provider);


  const deployResult = await hre.deployments.deploy("BSC_BscBridge", {
    contract: "BSC_BscBridge",
    ...await options(hre, tokenPairs,
      {
        sideBridgeAddress: sideBridgeAddress,
        wrappingTokenAddress: configFile.tokens.WBNB?.addresses.eth || ethers.constants.AddressZero,
        fee: 10,    // todo
        timeframeSeconds: isMainNet ? 60 * 60 * 4 : 60,
        lockTime: isMainNet ? 60 * 10 : 60,
        minSafetyBlocks: isMainNet ? 10 : 2,
      },
      [
        initialValidators,
        vsAddress,
      ])
  });


  configFile.bridges.eth.side = deployResult.address;
  configFile.save()

  if (deployResult.newlyDeployed) {
    console.log('Call this cmd second time to update tokens')
    return;
  }

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, "BSC_BscBridge");
};


async function getValidators(ambProvider: EthereumProvider): Promise<string[]> {
  const provider = new ethers.providers.JsonRpcProvider(urlFromHHProvider(ambProvider))

  const vsContract = ethers.ContractFactory.getContract(vsAddress, vsAbi)
  const block = await provider.getBlock('latest');  // todo block where Transfer event with eventId 0 emitted
  return await vsContract.connect(provider).getValidators({blockTag: block.number});
}


export default func;
func.tags = ["bridges_bnb"];
