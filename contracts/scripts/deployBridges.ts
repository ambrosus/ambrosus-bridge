import {ethers} from "hardhat";
import {Contract, Signer} from "ethers";
import vsAbi from "../abi/ModifiedValidatorSet.json"


interface CommonArgs {
  sideBridgeAddress: string | null;
  relayAddress: string;

  tokens: { [th: string]: any };

  fee: number;
  feeRecipient: string;

  timeframe: number;
  lockTime: number;
  minSafetyBlocks: number;

  ambWrapperAddress: string;
}


export async function deployAmbBridge(
  ambSigner: Signer,
  args: CommonArgs
): Promise<Contract> {
  const argsForAmb = [[
    ethers.constants.AddressZero,  // ethBridge not deployed yet
    args.relayAddress,

    Object.keys(args.tokens), Object.values(args.tokens),

    args.fee,
    args.feeRecipient,

    args.timeframe,
    args.lockTime,
    args.minSafetyBlocks,
  ],
  args.ambWrapperAddress
  ];

  return await deploy("AmbBridge", ambSigner, ...argsForAmb);
}


export async function deployEthBridge(
  ethSigner: Signer,
  args: CommonArgs,

  vsContractAddress: string,
  ambSigner: Signer,
): Promise<Contract> {
  const [initialValidators, lastProcessedBlock] = await getValidators(ambSigner, vsContractAddress);

  const argsForEth = [[
    args.sideBridgeAddress,
    args.relayAddress,

    Object.keys(args.tokens), Object.values(args.tokens),

    args.fee,
    args.feeRecipient,

    args.timeframe,
    args.lockTime,
    args.minSafetyBlocks,
  ],
    initialValidators,
    vsContractAddress,
    lastProcessedBlock
  ];

  return await deploy("EthBridge", ethSigner, ...argsForEth);
}


async function getValidators(ambSigner: Signer, vsContractAddress: string): Promise<[string[], string]> {
  const vsContract = ethers.ContractFactory.getContract(vsContractAddress, vsAbi, ambSigner)
  const block = await ambSigner.provider!.getBlock('latest');
  const validators = await vsContract["getValidators()"]({blockTag: block.number});
  return [validators, block.hash];
}


export async function deploy(contractName: string, signer: Signer, ...args: Array<any>) {
  const factory = await ethers.getContractFactory(contractName, signer);
  const deployedContract = await factory.deploy(...args);
  await deployedContract.deployed();
  console.log(`${contractName} deployed at ${deployedContract.address}; 
    tx: ${deployedContract.deployTransaction.hash}`);
  return deployedContract;
}
