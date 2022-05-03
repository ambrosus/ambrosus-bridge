import {EthereumProvider, HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import vsAbi from "../abi/ValidatorSet.json";
import {ethers} from "ethers";
import {
  addNewTokensToBridge,
  configPath,
  getTokenPairs,
  networkType,
  readConfig,
  urlFromHHProvider,
  writeConfig
} from "./utils";

const vsAddress = "0x0000000000000000000000000000000000000F00" // todo get from something?


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["bsc"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  const path = configPath(hre.network);
  let configFile = readConfig(path);

  const {owner, proxyAdmin} = await hre.getNamedAccounts();
  // todo get admin and relay from getNamedAccounts
  const admin = owner;
  const relay = owner;

  const tokenPairs = getTokenPairs("bsc", "amb", hre.network)

  const ambNet = hre.companionNetworks['amb']
  const {address: sideBridgeAddress} = await ambNet.deployments.get('BSC_AmbBridge');
  const initialValidators = await getValidators(ambNet.provider);


  const deployResult = await hre.deployments.deploy("BSC_BscBridge", {
    contract: "BSC_BscBridge",
    from: owner,
    proxy: {
      owner: proxyAdmin,
      proxyContract: "proxyTransparent",
      execute: {
        init: {
          methodName: "initialize",
          args: [
              {
              sideBridgeAddress: sideBridgeAddress,
              adminAddress: admin,
              relayAddress: relay,
              wrappingTokenAddress: configFile.tokens.WBNB?.addresses.eth || ethers.constants.AddressZero,
              tokenThisAddresses: Object.keys(tokenPairs),
              tokenSideAddresses: Object.values(tokenPairs),
              fee: 10,    // todo
              feeRecipient: owner,   // todo
              timeframeSeconds: isMainNet ? 60*60*4 : 60,
              lockTime: isMainNet ? 60*10 : 60,
              minSafetyBlocks: isMainNet ? 10 : 2,
            },
            initialValidators,
            vsAddress,
          ]
        }
      }
    },
    log: true
  });


  configFile.bridges.eth.side = deployResult.address;
  writeConfig(path, configFile);

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
