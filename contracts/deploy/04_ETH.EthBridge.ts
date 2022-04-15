import {EthereumProvider, HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import vsAbi from "../abi/ModifiedValidatorSet.json";
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
  if (hre.network.live && !hre.network.tags["eth"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  const path = configPath(hre.network);
  let configFile = readConfig(path);

  const {owner} = await hre.getNamedAccounts();
  // todo get admin and relay from getNamedAccounts
  const admin = owner;
  const relay = owner;

  const tokenPairs = getTokenPairs("eth", "amb", hre.network)

  const ambNet = hre.companionNetworks['amb']
  const {address: sideBridgeAddress} = await ambNet.deployments.get('AmbBridge');
  const [initialValidators, lastProcessedBlock] = await getValidators(ambNet.provider);


  const deployResult = await hre.deployments.deploy("EthBridge", {
    contract: "EthBridge",
    from: owner,
    args: [
      {
        sideBridgeAddress: sideBridgeAddress,
        adminAddress: admin,
        relayAddress: relay,
        wrappingTokenAddress: configFile.tokens.WETH?.addresses.eth || ethers.constants.AddressZero,
        tokenThisAddresses: Object.keys(tokenPairs),
        tokenSideAddresses: Object.values(tokenPairs),
        fee: 10,    // todo
        feeRecipient: owner,   // todo
        timeframeSeconds: isMainNet ? 14400 : 1,
        lockTime: isMainNet ? 1000 : 1,
        minSafetyBlocks: 10,
      },
      initialValidators,
      vsAddress,
      lastProcessedBlock,
    ],
    log: true,
    skipIfAlreadyDeployed: true
  });


  configFile.bridges.eth.side = deployResult.address;
  writeConfig(path, configFile);

  if (deployResult.newlyDeployed) {
    console.log('Call this cmd second time to update tokens')
    return;
  }

  // add new tokens
  await addNewTokensToBridge(tokenPairs, hre, "EthBridge");

};


async function getValidators(ambProvider: EthereumProvider): Promise<[string[], string]> {
  const provider = new ethers.providers.JsonRpcProvider(urlFromHHProvider(ambProvider))

  const vsContract = ethers.ContractFactory.getContract(vsAddress, vsAbi)
  const block = await provider.getBlock('latest');
  const validators = await vsContract.connect(provider).getValidators({blockTag: block.number});

  return [validators, block.hash];
}


export default func;
func.tags = ["bridges"];
