import {EthereumProvider, HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import vsAbi from "../abi/ModifiedValidatorSet.json";
import {ethers} from "ethers";
import {getTokensPair, urlFromHHProvider} from "./utils";

const vsAddress = "0x0000000000000000000000000000000000000F00" // todo get from something?
const relayAddress = "0x295c2707319ad4beca6b5bb4086617fd6f240cfe" // todo get from something?


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["eth"]) return;

  const {owner} = await hre.getNamedAccounts();

  const ambNet = hre.companionNetworks['amb']
  const {address: sideBridgeAddress} = await ambNet.deployments.get('AmbBridge');
  const [initialValidators, lastProcessedBlock] = await getValidators(ambNet.provider, vsAddress);

  const [tokensThis, tokensSide] = getTokensPair("eth", "amb", hre.network)


  await hre.deployments.deploy("EthBridge", {
    contract: "EthBridge",
    from: owner,
    args: [
      {
        sideBridgeAddress: sideBridgeAddress,
        relayAddress: relayAddress,
        tokenThisAddresses: tokensThis,
        tokenSideAddresses: tokensSide,
        fee: 10,    // todo
        feeRecipient: owner,   // todo
        timeframeSeconds: hre.network.live ? 14400 : 1,
        lockTime: hre.network.live ? 1000 : 1,
        minSafetyBlocks: 10,
      },
      initialValidators,
      vsAddress,
      lastProcessedBlock,
    ],
    log: true,
  });


};


async function getValidators(ambProvider: EthereumProvider, vsContractAddress: string): Promise<[string[], string]> {
  const provider = new ethers.providers.JsonRpcProvider(urlFromHHProvider(ambProvider))

  const vsContract = ethers.ContractFactory.getContract(vsContractAddress, vsAbi)
  const block = await provider.getBlock('latest');
  const validators = await vsContract.connect(provider).getValidators({blockTag: block.number});
  return [validators, block.hash];
}


export default func;
func.tags = ["ethbridge"];
func.dependencies = ['AmbBridge', 'BridgeERC20Tokens'];
