import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.name !== 'hardhat' && hre.network.name !== 'amb') return
  const {owner} = await hre.getNamedAccounts();
  await hre.deployments.deploy("ValidatorSet", {
    from: owner,
    args: [
      ["0x195c2707319ad4beca6b5bb4086617fd6f240cfe", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe", "0x395c2707319ad4beca6b5bb4086617fd6f240cfe"], // validators
      1000, // block number
    ],
    log: true,
  });
};



export default func;
func.tags = ["validatorset"];
