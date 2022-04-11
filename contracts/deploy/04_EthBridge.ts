import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";
import vsAbi from "../abi/ModifiedValidatorSet.json";
import { ethers, VoidSigner } from "ethers";

const vsAddress = "0x0000000000000000000000000000000000000F00" // todo get from something?
const relayAddress = "0x295c2707319ad4beca6b5bb4086617fd6f240cfe" // todo get from something?


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["eth"]) return;

  const { owner } = await hre.getNamedAccounts();

  const ambNet = hre.companionNetworks['amb']
  const {address: sideBridgeAddress} = await ambNet.deployments.get('AmbBridge');
  const [initialValidators, lastProcessedBlock] = await getValidators(ambNet.provider as any, vsAddress);


  await hre.deployments.deploy("EthBridge", {
    contract: "EthBridge",
    from: owner,
    args: [
      {
        sideBridgeAddress: sideBridgeAddress,
        relayAddress: relayAddress,
        tokenThisAddresses: [],
        tokenSideAddresses: [],
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

  // todo set sideBridgeAddress in amb bridge

};


async function getValidators(ambProvider: ethers.providers.Provider, vsContractAddress: string): Promise<[string[], string]> {
  const ambSigner = new VoidSigner(ethers.constants.AddressZero, ambProvider);
  const vsContract = ethers.ContractFactory.getContract(vsContractAddress, vsAbi, ambSigner)
  const block = await ambSigner.provider!.getBlock('latest');
  const validators = await vsContract["getValidators()"]({blockTag: block.number});
  return [validators, block.hash];
}


export default func;
func.dependencies = ['AmbBridge'];
func.tags = ["ethbridge"];
