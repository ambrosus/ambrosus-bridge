import {EthereumProvider, HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import vsAbi from "../abi/ModifiedValidatorSet.json";
import {ethers} from "ethers";
import {configPath, getTokensPair, readConfig, urlFromHHProvider, writeConfig} from "./utils";

const vsAddress = "0x0000000000000000000000000000000000000F00" // todo get from something?
const relayAddress = "0x295c2707319ad4beca6b5bb4086617fd6f240cfe" // todo get from something?


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["eth"]) {
    console.log("Skipping ETH bridge deployment on non-ETH network")
    return;
  }

  const {owner} = await hre.getNamedAccounts();

  const ambNet = hre.companionNetworks['amb']
  const {address: sideBridgeAddress} = await ambNet.deployments.get('AmbBridge');
  const [initialValidators, lastProcessedBlock] = await getValidators(ambNet.provider, hre);


  const [tokensThis, tokensSide] = getTokensPair("eth", "amb", hre.network)

  const deployResult = await hre.deployments.deploy("EthBridge", {
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


  const path = configPath(hre.network);
  let configFile = readConfig(path);
  configFile.bridges.eth.side = deployResult.address;
  writeConfig(path, configFile);
};


async function getValidators(ambProvider: EthereumProvider, hre: HardhatRuntimeEnvironment): Promise<[string[], string]> {
  const provider = new ethers.providers.JsonRpcProvider(urlFromHHProvider(ambProvider))

  const vsContract = ethers.ContractFactory.getContract(vsAddress, vsAbi)
  const block = await provider.getBlock('latest');
  const validators = await vsContract.connect(provider).getValidators({blockTag: block.number});

  // contract will be redeployed if args are changed
  // check that validatorSet is not changed, then use lastProcessedBlock from prev deployment
  const oldDeployment = await hre.deployments.getOrNull("EthBridge");
  if (oldDeployment && oldDeployment.args) {
    const [_0, oldValidators, _2, oldLastProcessedBlock] = oldDeployment.args;
    if (JSON.stringify(validators) === JSON.stringify(oldValidators))  // i'm fucking hate js
      return [validators, oldLastProcessedBlock];
  }

  return [validators, block.hash];
}


export default func;
func.tags = ["bridges"];
