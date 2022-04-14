import {EthereumProvider, HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import vsAbi from "../abi/ModifiedValidatorSet.json";
import {ethers} from "ethers";
import {
  addNewTokensToBridge,
  configPath,
  getTokenPairs,
  networkType,
  readConfig, setAdminRole,
  urlFromHHProvider,
  writeConfig
} from "./utils";

const vsAddress = "0x0000000000000000000000000000000000000F00" // todo get from something?
const relayAddress = "0x295c2707319ad4beca6b5bb4086617fd6f240cfe" // todo get from something?


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["eth"]) return;
  const isMainNet = networkType(hre.network) === 'mainnet'

  const ambNet = hre.companionNetworks['amb']
  const {address: sideBridgeAddress} = await ambNet.deployments.get('AmbBridge');
  const [initialValidators, lastProcessedBlock] = await getValidators(ambNet.provider);

  const {owner} = await hre.getNamedAccounts();
  const tokenPairs = getTokenPairs("amb", "eth", hre.network)

  const deployResult = await hre.deployments.deploy("EthBridge", {
    contract: "EthBridge",
    from: owner,
    args: [
      {
        sideBridgeAddress: sideBridgeAddress,
        relayAddress: relayAddress,
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


  const path = configPath(hre.network);
  let configFile = readConfig(path);
  configFile.bridges.eth.side = deployResult.address;
  writeConfig(path, configFile);

  // set adminRole
  await setAdminRole(hre, "EthBridge")

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
